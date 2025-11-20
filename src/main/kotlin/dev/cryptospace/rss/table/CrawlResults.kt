package dev.cryptospace.rss.table

import org.jetbrains.exposed.dao.id.UUIDTable
import org.jetbrains.exposed.sql.ReferenceOption

object CrawlResults : UUIDTable(name = "crawl_results") {
    val target =
        reference(
            name = "target_id",
            refColumn = CrawlTargets.id,
            onDelete = ReferenceOption.CASCADE,
        ).index()
    val body = blob(name = "body")
}
