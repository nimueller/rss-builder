package dev.cryptospace.rss

import dev.cryptospace.rss.entity.CrawlTarget
import kotlinx.coroutines.delay
import org.openqa.selenium.By
import org.openqa.selenium.WebElement
import org.openqa.selenium.firefox.FirefoxDriver
import org.openqa.selenium.firefox.FirefoxOptions
import kotlin.time.Duration.Companion.milliseconds

private const val GECKO_DRIVER_PROPERTY = "webdriver.gecko.driver"

object Crawler {
    init {
        val driverPath = Crawler.javaClass.getResource("/geckodriver")?.path

        if (driverPath != null) {
            System.setProperty(GECKO_DRIVER_PROPERTY, driverPath)
        }
    }

    private val webDriver = FirefoxDriver(FirefoxOptions().setHeadless(true))

    suspend fun CrawlTarget.open() {
        // just calling webdriver[url] would be ugly as the getter returns void
        @Suppress("kotlin:S6518")
        webDriver.get(url)

        adBannerWaitTimeInMillis?.let { delay(it.milliseconds) }
        adBannerButtonSelector?.let { webDriver.findElement(By.cssSelector(it)).click() }
    }

    fun CrawlTarget.fetchItems(): List<WebElement> = webDriver.findElements(By.cssSelector(itemSelector))
}
