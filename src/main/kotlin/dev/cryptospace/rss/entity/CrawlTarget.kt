package dev.cryptospace.rss.entity

import dev.cryptospace.rss.table.CrawlTargets
import org.jetbrains.exposed.dao.UUIDEntity
import org.jetbrains.exposed.dao.UUIDEntityClass
import org.jetbrains.exposed.dao.id.EntityID
import java.util.UUID

class CrawlTarget(
    id: EntityID<UUID>,
) : UUIDEntity(id) {
    var url by CrawlTargets.url
    var adBannerWaitTimeInMillis by CrawlTargets.adBannerWaitTimeInMillis
    var adBannerButtonSelector by CrawlTargets.adBannerButtonSelector
    var itemSelector by CrawlTargets.itemSelector
    var itemTitleXPath by CrawlTargets.itemTitleXPath
    var itemLinkXPath by CrawlTargets.itemLinkXPath

    companion object : UUIDEntityClass<CrawlTarget>(CrawlTargets)
}
