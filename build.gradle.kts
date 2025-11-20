plugins {
    application
    alias(libs.plugins.kotlin.jvm)
    alias(libs.plugins.ktlint)
    alias(libs.plugins.sonarlint)
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
}

tasks.test {
    useJUnitPlatform()
}
