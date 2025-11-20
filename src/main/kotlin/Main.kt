import dev.cryptospace.rss.Crawler.fetchItems
import dev.cryptospace.rss.Crawler.open
import dev.cryptospace.rss.entity.CrawlTarget
import dev.cryptospace.rss.table.CrawlResults
import dev.cryptospace.rss.table.CrawlTargets
import kotlinx.coroutines.launch
import kotlinx.coroutines.runBlocking
import org.jetbrains.exposed.sql.Database
import org.jetbrains.exposed.sql.SchemaUtils
import org.jetbrains.exposed.sql.transactions.transaction

fun main() {
    runBlocking {
        val dbHost = System.getenv("DB_HOST") ?: "localhost"
        val dbPort = (System.getenv("DB_PORT") ?: "5432").toInt()
        val dbName = System.getenv("POSTGRES_DB") ?: System.getenv("DB_NAME") ?: "rss-builder"
        val dbUser = System.getenv("POSTGRES_USER") ?: System.getenv("DB_USER") ?: "postgres"
        val dbPass = System.getenv("POSTGRES_PASSWORD") ?: System.getenv("DB_PASSWORD") ?: "postgres"

        val jdbcUrl = "jdbc:postgresql://$dbHost:$dbPort/$dbName"
        Database.connect(jdbcUrl, driver = "org.postgresql.Driver", user = dbUser, password = dbPass)

        transaction {
            SchemaUtils.createMissingTablesAndColumns(CrawlTargets, CrawlResults)

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
