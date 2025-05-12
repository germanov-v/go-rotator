
CREATE TABLE slots (
                       id   TEXT PRIMARY KEY
);

CREATE TABLE banners (
                         id   TEXT PRIMARY KEY
);

CREATE TABLE groups (
                        id   TEXT PRIMARY KEY
);

CREATE TABLE slot_banners (
                              slot_id   TEXT NOT NULL
                                  REFERENCES slots(id)
                                      ON DELETE CASCADE,
                              banner_id TEXT NOT NULL
                                  REFERENCES banners(id)
                                      ON DELETE CASCADE,
                              PRIMARY KEY (slot_id, banner_id)
);

CREATE TABLE stats (
                       slot_id     TEXT NOT NULL
                           REFERENCES slots(id)
                               ON DELETE CASCADE,
                       banner_id   TEXT NOT NULL
                           REFERENCES banners(id)
                               ON DELETE CASCADE,
                       group_id    TEXT NOT NULL
                           REFERENCES groups(id)
                               ON DELETE CASCADE,
                       impressions BIGINT NOT NULL DEFAULT 0,
                       clicks      BIGINT NOT NULL DEFAULT 0,
                       PRIMARY KEY (slot_id, banner_id, group_id)
);
