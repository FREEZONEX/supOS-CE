---
name: UNS Validation
description: Validate UNS namespace JSON structures against node type rules, topicType hierarchy, field type constraints, and naming conventions.
---

# UNS Validation Skill

This skill provides validation rules and procedures for checking UNS (Unified Namespace) namespace JSON structures.

---

## Quick Validation Summary

| Check | Rule |
|-------|------|
| Node type | Every node has `type` = `PATH` or `TOPIC` |
| Non-empty name | Every node has a non-empty `name` |
| Type folder | Every TOPIC is a descendant of a PATH with `topicType` |
| topicType values | Only `STATE`, `ACTION`, or `METRIC` |
| TOPIC fields | Has `fields` array or `templateAlias` (at least one) |
| METRIC fields | Only numeric (`INTEGER`, `LONG`, `FLOAT`, `DOUBLE`) and `DATETIME` types |
| Boolean strings | `enableHistory`, `mockData`, `writeData` are `"TRUE"` / `"FALSE"` strings |
| Folder naming | Not reserved (`label`, `template`), max 63 chars |

---

## Validation Rules

### Rule V1: Valid Node Type

Every node must have `type` set to `PATH` or `TOPIC` (case-insensitive match, conventionally uppercase).

```
ASSERT: node.type IN ["PATH", "TOPIC"]
```

### Rule V2: Non-Empty Name

Every node must have a non-empty `name`.

```
ASSERT: node.name != "" AND node.name != null
```

### Rule V3: TOPIC Under Type Folder

Every TOPIC must be a descendant of a PATH node that has `topicType` set. It does not need to be a direct child — the `topicType` propagates through intermediate PATH nodes.

```
ASSERT: FOR EACH topic node,
        EXISTS ancestor WHERE ancestor.type == "PATH" AND ancestor.topicType IN ["STATE", "ACTION", "METRIC"]
```

```
# Valid: TOPIC is direct child of type folder
v1 / Plant / Area / State(topicType:STATE) / current_job(TOPIC)  ✓

# Valid: TOPIC is nested under type folder via intermediate PATH
v1 / Plant / Area / State(topicType:STATE) / SubGroup(PATH) / current_job(TOPIC)  ✓

# Invalid: TOPIC has no type folder ancestor
v1 / Plant / Area / current_job(TOPIC)  ✗
```

### Rule V4: topicType Values

When a PATH node has `topicType`, it must be one of: `STATE`, `ACTION`, `METRIC`.

```
ASSERT: IF node.topicType IS SET THEN node.topicType IN ["STATE", "ACTION", "METRIC"]
```

### Rule V5: TOPIC Has Fields or Template

A TOPIC must define its schema via `fields` array or `templateAlias` (at least one).

```
ASSERT: IF node.type == "TOPIC" THEN (len(node.fields) > 0 OR node.templateAlias != "")
```

### Rule V6: METRIC Field Types

TOPICs under a `METRIC` type folder should only use numeric and datetime field types. String/Boolean fields in METRIC topics are likely errors.

```
ASSERT: IF inherited topicType == "METRIC" THEN
        FOR EACH field IN node.fields:
            field.type IN ["INTEGER", "LONG", "FLOAT", "DOUBLE", "DATETIME"]
```

### Rule V7: Valid Field Types

All field types must be one of the supported types.

```
ASSERT: field.type IN ["INTEGER", "LONG", "FLOAT", "DOUBLE", "BOOLEAN", "DATETIME", "STRING"]
```

> Field type matching is case-insensitive. `"int"` is accepted as alias for `INTEGER`.

### Rule V8: Boolean String Fields

`enableHistory`, `mockData`, and `writeData` must be string values `"TRUE"` or `"FALSE"`, not JSON booleans.

```
ASSERT: IF node.enableHistory IS SET THEN node.enableHistory IN ["TRUE", "FALSE"]
ASSERT: IF node.mockData IS SET THEN node.mockData IN ["TRUE", "FALSE"]
ASSERT: IF node.writeData IS SET THEN node.writeData IN ["TRUE", "FALSE"]
```

### Rule V9: Folder Name Constraints

PATH node names must not use reserved words and must respect length limits.

```
ASSERT: IF node.type == "PATH" THEN node.name NOT IN ["label", "template"]
ASSERT: IF node.type == "PATH" THEN len(node.name) <= 63
```

### Rule V10: Template Reference Exists

If a TOPIC uses `templateAlias`, the referenced template must be defined in the `templates` array.

```
ASSERT: IF node.templateAlias IS SET THEN
        EXISTS template IN templates WHERE template.alias == node.templateAlias
```

### Rule V11: extendProperties Is Flat Map

`extendProperties` must be a flat key-value object (not nested).

```
ASSERT: IF node.extendProperties IS SET THEN typeof(node.extendProperties) == "object"
```

---

## Common Validation Errors

| Error | Cause | Fix |
|-------|-------|-----|
| Missing type folder | TOPIC not under any PATH with `topicType` | Add a parent PATH with `topicType: STATE/ACTION/METRIC` |
| Invalid field type | Field uses unsupported type string | Use one of: `INTEGER`, `LONG`, `FLOAT`, `DOUBLE`, `BOOLEAN`, `DATETIME`, `STRING` |
| Non-numeric METRIC field | METRIC topic has `STRING` or `BOOLEAN` field | Move to STATE, or change field types to numeric |
| Reserved folder name | PATH named `label` or `template` | Rename the folder |
| Boolean value not string | `enableHistory: true` instead of `"TRUE"` | Use string `"TRUE"` or `"FALSE"` |
| Missing fields | TOPIC has no `fields` and no `templateAlias` | Add `fields` array or set `templateAlias` |
| Undefined template | `templateAlias` references non-existent template | Define template in `templates` array or fix alias |

---

## Validation Pseudocode

```python
def validate_uns_tree(node, inherited_topic_type=None, templates=None):
    errors = []
    templates = templates or {}

    # V1: Valid node type
    if node["type"].upper() not in ("PATH", "TOPIC"):
        errors.append(f"Invalid type '{node['type']}' on node '{node.get('name')}'")
        return errors

    # V2: Non-empty name
    if not node.get("name"):
        errors.append("Node has empty or missing name")

    is_path = node["type"].upper() == "PATH"
    is_topic = node["type"].upper() == "TOPIC"

    # Determine effective topicType (own or inherited)
    current_topic_type = node.get("topicType", inherited_topic_type)

    if is_path:
        # V4: Valid topicType
        if "topicType" in node and node["topicType"] not in ("STATE", "ACTION", "METRIC"):
            errors.append(f"Invalid topicType '{node['topicType']}' on PATH '{node['name']}'")

        # V9: Folder name constraints
        if node["name"] in ("label", "template"):
            errors.append(f"Reserved folder name: '{node['name']}'")
        if len(node["name"]) > 63:
            errors.append(f"Folder name exceeds 63 chars: '{node['name']}'")

        # Recurse children
        for child in node.get("children", []):
            errors.extend(validate_uns_tree(child, current_topic_type, templates))

    elif is_topic:
        # V3: Must be under a type folder
        if current_topic_type not in ("STATE", "ACTION", "METRIC"):
            errors.append(f"TOPIC '{node['name']}' not under a type folder (no inherited topicType)")

        # V5: Must have fields or templateAlias
        has_fields = len(node.get("fields", [])) > 0
        has_template = bool(node.get("templateAlias"))
        if not has_fields and not has_template:
            errors.append(f"TOPIC '{node['name']}' has no fields and no templateAlias")

        # V10: Template reference exists
        if has_template and node["templateAlias"] not in templates:
            errors.append(f"TOPIC '{node['name']}' references undefined template '{node['templateAlias']}'")

        # V7 + V6: Field type validation
        for field in node.get("fields", []):
            ft = field.get("type", "").upper()
            if ft not in ("INTEGER", "LONG", "FLOAT", "DOUBLE", "BOOLEAN", "DATETIME", "STRING"):
                errors.append(f"Invalid field type '{field['type']}' in TOPIC '{node['name']}'")
            elif current_topic_type == "METRIC" and ft in ("STRING", "BOOLEAN"):
                errors.append(f"METRIC TOPIC '{node['name']}' has non-numeric field '{field['name']}' ({ft})")

        # V8: Boolean string fields
        for flag in ("enableHistory", "mockData", "writeData"):
            if flag in node and node[flag] not in ("TRUE", "FALSE"):
                errors.append(f"TOPIC '{node['name']}': {flag} must be 'TRUE' or 'FALSE', got '{node[flag]}'")

    return errors


def validate_import(data):
    errors = []
    # Collect template aliases
    templates = {}
    for t in data.get("templates", []):
        if t.get("alias"):
            templates[t["alias"]] = t

    # Validate each root namespace node
    for root in data.get("namespace", []):
        errors.extend(validate_uns_tree(root, templates=templates))

    return errors
```

---

## Integration

When implementing changes to UNS structures:

1. **Before import**: Validate JSON structure against all rules above
2. **On review**: Focus on V3 (type folder ancestry) and V6 (METRIC field types) — the most common issues
3. **Auto-fix candidates**: Missing `topicType` on obvious type folders (named `State`/`Action`/`Metric`)

Refer to [UNS Structure Design](../uns-structure/SKILL.md) for the complete design guidelines.
