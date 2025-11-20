package dev.cryptospace.rss.entity

import dev.cryptospace.rss.table.CrawlResults
import org.jetbrains.exposed.dao.UUIDEntity
import org.jetbrains.exposed.dao.UUIDEntityClass
import org.jetbrains.exposed.dao.id.EntityID
import java.util.UUID

class CrawlResult(
    id: EntityID<UUID>,
) : UUIDEntity(id) {
    var target by CrawlTarget referencedOn CrawlResults.target
    var body by CrawlResults.body

    companion object : UUIDEntityClass<CrawlResult>(CrawlResults)
}
