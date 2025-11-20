import dev.cryptospace.rss.Crawler.fetchItems
import dev.cryptospace.rss.Crawler.open
import dev.cryptospace.rss.entity.CrawlTarget
import kotlinx.coroutines.launch
import kotlinx.coroutines.runBlocking
import liquibase.Liquibase
import liquibase.database.DatabaseFactory
import liquibase.resource.ClassLoaderResourceAccessor
import org.jetbrains.exposed.sql.Database
import org.jetbrains.exposed.sql.transactions.transaction

fun main() {
    runBlocking {
        val dbHost = System.getenv("DB_HOST") ?: "localhost"
        val dbPort = (System.getenv("DB_PORT") ?: "5432").toInt()
        val dbName = System.getenv("POSTGRES_DB") ?: System.getenv("DB_NAME") ?: "rss-builder"
        val dbUser = System.getenv("POSTGRES_USER") ?: System.getenv("DB_USER") ?: "postgres"
        val dbPass = System.getenv("POSTGRES_PASSWORD") ?: System.getenv("DB_PASSWORD") ?: "postgres"

        val jdbcUrl = "jdbc:postgresql://$dbHost:$dbPort/$dbName"

        // Now connect Exposed
        Database.connect(jdbcUrl, driver = "org.postgresql.Driver", user = dbUser, password = dbPass)

        // Seed initial data if empty
        transaction {
            if (CrawlTarget.all().empty()) {
                CrawlTarget.new {
                    url = "https://kicker.de"
                    adBannerWaitTimeInMillis = 1000
                    adBannerButtonSelector = "a.kick__btn-primary:nth-child(3)"
                    itemSelector = "kick__slick-slide"
                    itemTitleXPath = "/div/div/div/div/h3/a"
                    itemLinkXPath = "/div/div/div/div/h3/a@href"
                }
            }
        }

        val targets =
            transaction {
                CrawlTarget.all().toList()
            }

        targets.forEach { target ->
            launch {
                target.open()
                val items = target.fetchItems()
                println("Found ${items.size} items")
            }
        }
    }
}
