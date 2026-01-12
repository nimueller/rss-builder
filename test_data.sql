INSERT INTO scrap_target (url, base_goquery_selector, item_goquery_selector, image_goquery_selector,
                          article_goquery_selector)
VALUES ('https://kicker.de',
        '#kick__ressort > section:nth-child(2)',
        '.kick__slidelist__item',
        '.kick__slidelist__item_content_picture img',
        'main.kick__article__content');