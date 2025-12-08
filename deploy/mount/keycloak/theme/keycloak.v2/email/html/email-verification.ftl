<#import "template-custom.ftl" as layout>
<@layout.customEmailLayout>
${kcSanitize(msg("emailVerificationBodyHtml",link, linkExpiration, realmName, linkExpirationFormatter(linkExpiration)))?no_esc}
</@layout.customEmailLayout>
