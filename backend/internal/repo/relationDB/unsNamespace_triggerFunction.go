package relationDB

import (
	"log"

	"gorm.io/gorm"
)

func (p UnsNamespaceRepo) migrate(db *gorm.DB) {
	err := createFunction_nextId(db)
	if err != nil {
		log.Println("ERR: createFunction_nextId", err.Error())
	}
	err = createFunction_nextIdLong(db)
	if err != nil {
		log.Println("ERR: createFunction_nextIdLong", err.Error())
	}
	err = createFunction_extractFieldsText(db)
	if err != nil {
		log.Println("ERR: createFunction_extractFieldsText", err.Error())
	}
	err = createFunc_pathHash(db)
	if err != nil {
		log.Println("ERR: createFunc_pathHash", err.Error())
	}
	err = createTrigger_extractFieldsText(db)
	if err == nil {
		refreshFieldsText(db)
		log.Println("刷新 FieldsText")
	}
	err = createTrigger_pathHash(db)
	if err == nil {
		refreshPathHash(db)
		log.Println("刷新 PathHash")
	}
}
func createFunction_nextIdLong(db *gorm.DB) error {
	return db.Exec(`CREATE OR REPLACE FUNCTION "nextIdLong"(prev TEXT, layRec TEXT)
        RETURNS int8 AS $$
        DECLARE
        result TEXT;
        BEGIN
        result := "nextId"(prev, layRec);

        IF result = '' THEN
        RETURN -1;
        ELSE
        RETURN result::int8;
        END IF;
        END;
        $$ LANGUAGE plpgsql IMMUTABLE`).Error
}
func createFunction_nextId(db *gorm.DB) error {
	return db.Exec(`
        CREATE OR REPLACE FUNCTION "nextId" (prev TEXT, layRec TEXT)
        RETURNS TEXT AS $$
        DECLARE
        prev_parts TEXT[];
        rec_parts TEXT[];
        i INT;
        match BOOLEAN := TRUE;
        BEGIN
        -- 如果 prev 为空，直接返回 layRec 的第一部分
        IF prev = '' THEN
        RETURN split_part(layRec, '/', 1);
        END IF;

        -- 将字符串分割为数组
        prev_parts := string_to_array(prev, '/');
        rec_parts := string_to_array(layRec, '/');

        -- 检查 prev 是否与 layRec 的前面部分匹配
        FOR i IN 1..array_length(prev_parts, 1) LOOP
        IF i > array_length(rec_parts, 1) OR prev_parts[i] != rec_parts[i] THEN
        match := FALSE;
        EXIT;
        END IF;
        END LOOP;

        -- 如果不匹配，返回空字符串
        IF NOT match THEN
        RETURN '';
        END IF;

        -- 如果完全匹配，检查是否有下一级
        IF array_length(prev_parts, 1) = array_length(rec_parts, 1) THEN
        RETURN '';
        ELSE
        -- 返回下一级
        RETURN rec_parts[array_length(prev_parts, 1) + 1];
        END IF;
        END;
        $$ LANGUAGE plpgsql IMMUTABLE`).Error
}
func createFunction_extractFieldsText(db *gorm.DB) error {
	return db.Exec(`
        CREATE OR REPLACE FUNCTION extract_fields_text()
        RETURNS TRIGGER AS $$
        DECLARE
        field_item JSONB;
        extracted_text TEXT := '';
        is_system_field BOOLEAN;
        BEGIN
        -- 确保 fields 字段不为空
        IF NEW.fields IS NOT NULL THEN
        -- 遍历 JSONB 数组中的每个对象
        FOR field_item IN SELECT * FROM jsonb_array_elements(NEW.fields::jsonb)
        LOOP
        -- 检查 systemField 是否为 false
        is_system_field := COALESCE((field_item->>'systemField')::BOOLEAN, false);

        -- 只有当 systemField 为 false 时才提取 name 和 remark
        IF NOT is_system_field THEN
        -- 提取每个对象的 name 和 remark 字段，用空格连接
        extracted_text := extracted_text || ' ' ||
        COALESCE(field_item->>'name', '') || ' ' ||
        COALESCE(field_item->>'remark', '');
        END IF;
        END LOOP;

        -- 去除首尾空格并更新 fields_text 字段
        NEW.fields_text := TRIM(extracted_text);
        ELSE
        NEW.fields_text := NULL;
        END IF;

        RETURN NEW;
        END;
        $$ LANGUAGE plpgsql
   `).Error
}
func createTrigger_extractFieldsText(db *gorm.DB) error {
	return db.Exec(`
        CREATE TRIGGER trigger_extract_fields_text
        BEFORE INSERT OR UPDATE OF fields ON uns_namespace
        FOR EACH ROW
        EXECUTE FUNCTION extract_fields_text(
    `).Error
}
func createFunc_pathHash(db *gorm.DB) error {
	return db.Exec(`
        CREATE OR REPLACE FUNCTION update_pathash()
        RETURNS TRIGGER AS $$
        BEGIN
        NEW.pathash := hashtext(NEW.path);
        RETURN NEW;
        END;
        $$ LANGUAGE plpgsql
    `).Error
}
func createTrigger_pathHash(db *gorm.DB) error {
	return db.Exec(`
        CREATE TRIGGER trigger_update_pathash
        BEFORE INSERT OR UPDATE OF path ON uns_namespace
        FOR EACH ROW
        EXECUTE FUNCTION update_pathash()
    `).Error
}
func refreshFieldsText(db *gorm.DB) error {
	return db.Exec(`UPDATE uns_namespace SET fields = fields WHERE fields IS NOT NULL`).Error
}
func refreshPathHash(db *gorm.DB) error {
	return db.Exec(`UPDATE uns_namespace SET path=path WHERE fields IS NOT NULL`).Error
}
