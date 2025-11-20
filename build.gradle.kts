buildscript {
    repositories {
        mavenCentral()
    }
    dependencies {
        // Liquibase itself plugin must be in the buildscript classpath for some strange reason
        // https://github.com/liquibase/liquibase-gradle-plugin/blob/master/doc/usage.md#2-setting-up-the-classpath
        classpath(libs.liquibase.core)
    }
}

plugins {
    application
    alias(libs.plugins.kotlin.jvm)
    alias(libs.plugins.ktlint)
    alias(libs.plugins.sonarlint)
    alias(libs.plugins.liquibase)
}

group = "dev.cryptospace"
version = "1.0-SNAPSHOT"

application {
    mainClass.set("MainKt")
}

repositories {
    mavenCentral()
}

dependencies {
    testImplementation(libs.kotlin.test)

    sonarlintPlugins(libs.sonar.kotlin.plugin)

    implementation(libs.postgresql)

    implementation(libs.exposed.core)
    implementation(libs.exposed.crypt)
    implementation(libs.exposed.dao)
    implementation(libs.exposed.jdbc)
    implementation(libs.exposed.kotlin.datetime)

    implementation(libs.selenium.java)
    implementation(libs.selenium.firefox.driver)

    // Database migrations (runtime inside app)
    implementation(libs.liquibase.core)

    // Liquibase Gradle plugin runtime (for ./gradlew update)
    liquibaseRuntime(libs.liquibase.core)
    liquibaseRuntime(libs.postgresql)
    liquibaseRuntime(libs.commons.lang)
    liquibaseRuntime(libs.picocli)
}

tasks.test {
    useJUnitPlatform()
}

// Configure Liquibase Gradle plugin to use env vars or sensible defaults
val dbHost = System.getenv("DB_HOST") ?: "localhost"
val dbPort = (System.getenv("DB_PORT") ?: "5432").toInt()
val dbName = System.getenv("POSTGRES_DB") ?: System.getenv("DB_NAME") ?: "rss-builder"
val dbUser = System.getenv("POSTGRES_USER") ?: System.getenv("DB_USER") ?: "postgres"
val dbPass = System.getenv("POSTGRES_PASSWORD") ?: System.getenv("DB_PASSWORD") ?: "postgres"

val jdbcUrl = "jdbc:postgresql://$dbHost:$dbPort/$dbName"

liquibase {
    activities.register("main") {
        arguments = mapOf(
            "logLevel" to "info",
            "changelogFile" to "src/main/resources/db/changelog/db.changelog-master.yaml",
            "url" to jdbcUrl,
            "username" to dbUser,
            "password" to dbPass,
        )
    }
    runList = "main"
}

// Ensure DB is migrated before launching the app with `./gradlew run`
tasks.named("run") {
    dependsOn("update")
}
