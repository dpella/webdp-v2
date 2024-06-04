CREATE TYPE WebDPType AS ENUM ('Int', 'Double', 'Bool', 'Text', 'Enum');
CREATE TYPE PrivacyNotion AS ENUM ('PureDP', 'ApproxDP');


CREATE TABLE Roles (
    role TEXT PRIMARY KEY
);

CREATE TABLE Users (
    handle TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    pwd TEXT NOT NULL,
    created_time TIMESTAMPTZ NOT NULL,
    updated_time TIMESTAMPTZ NOT NULL,
    CHECK(LENGTH(handle) > 0),
    CHECK(LENGTH(pwd) > 0),
    CHECK(LENGTH(name) > 0)
);

CREATE TABLE UserTokens (
    username TEXT,
    token TEXT NOT NULL,
    PRIMARY KEY (username),
    FOREIGN KEY (username) REFERENCES Users(handle) ON DELETE CASCADE
);

CREATE TABLE UserRoles (
    username TEXT,
    role TEXT,
    PRIMARY KEY (username, role),
    FOREIGN KEY (username) REFERENCES Users(handle) ON DELETE CASCADE,
    FOREIGN KEY (role) REFERENCES Roles(role) ON DELETE CASCADE
);

CREATE TABLE Dataset (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    owner TEXT, 
    privacy_notion PrivacyNotion NOT NULL,
    total_epsilon DOUBLE PRECISION NOT NULL, 
    total_delta DOUBLE PRECISION,
    created_time TIMESTAMPTZ NOT NULL,
    updated_time TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (owner) REFERENCES Users(handle) ON DELETE CASCADE,
    CHECK(privacy_notion != 'PureDP' OR (COALESCE(total_delta, 0.0) = 0.0)),
    CHECK(privacy_notion != 'ApproxDP' OR total_delta IS NOT NULL),
    CHECK(total_epsilon >= 0.0),
    CHECK(COALESCE(total_delta, 0.0) >= 0.0)
);

/*
,
    CHECK(privacy_notion != 'PureDP' OR (total_delta IS NULL or total_delta = 0.0)),
    CHECK(total_epsilon >= 0.0),
    CHECK(COALESCE(total_delta, 0.0) >= 0.0)

*/

CREATE TABLE DataUpload (
    dataset SERIAL PRIMARY KEY,
    loaded_time TIMESTAMPTZ NOT NULL,
    loaded_data BYTEA NOT NULL,
    FOREIGN KEY (dataset) REFERENCES Dataset(id) ON DELETE CASCADE
);

CREATE TABLE ColumnSchemas (
    dataset SERIAL, 
    column_name TEXT,
    data_type WebDPType NOT NULL, 
    low INT,
    high INT,
    labels TEXT[],
    PRIMARY KEY (dataset, column_name),
    FOREIGN KEY (dataset) REFERENCES Dataset(id) ON DELETE CASCADE,
    CHECK (
        (data_type = 'Int' OR data_type = 'Double') AND low IS NOT NULL AND high IS NOT NULL AND labels IS NULL OR 
        (data_type = 'Bool' OR data_type = 'Text') AND low IS NULL AND high IS NULL AND labels IS NULL OR 
        (data_type = 'Enum') AND low IS NULL AND high IS NULL and labels IS NOT NULL
        )
);

CREATE TABLE DPEngines (
    name TEXT PRIMARY KEY, 
    eval_url TEXT,
    accuracy_url TEXT 
);

CREATE TABLE UserBudgetAllocation (
    dataset SERIAL,
    userid TEXT,
    all_epsilon DOUBLE PRECISION NOT NULL,
    all_delta DOUBLE PRECISION,
    con_epsilon DOUBLE PRECISION,
    con_delta DOUBLE PRECISION,
    PRIMARY KEY (dataset, userid),
    FOREIGN KEY (dataset) REFERENCES Dataset(id) ON DELETE CASCADE,
    FOREIGN KEY (userid) REFERENCES Users(handle) ON DELETE CASCADE,
    CHECK (all_epsilon >= COALESCE(con_epsilon, 0)),
    CHECK (COALESCE(all_delta, 0) >= COALESCE(con_delta, 0))
);


INSERT INTO Roles VALUES ('Analyst');
INSERT INTO Roles VALUES ('Admin');
INSERT INTO Roles VALUES ('Curator');


CREATE VIEW LoadedDatasets AS (
    SELECT D.id, 
    D.name, 
    D.owner, 
    D.privacy_notion, 
    D.total_epsilon, 
    COALESCE(D.total_delta, 0) AS total_delta, 
    CASE WHEN L.loaded_time IS NULL THEN false ELSE true END AS loaded, 
    D.created_time,
    D.updated_time,
    L.loaded_time
    FROM Dataset AS D LEFT OUTER JOIN DataUpload AS L ON D.id = L.dataset
);

CREATE VIEW DatasetAllocatedConsumed AS (
    WITH Consumed AS (
        SELECT dataset, 
        SUM(all_epsilon) AS all_epsilon,
        SUM(all_delta) AS all_delta, 
        SUM(con_epsilon) AS con_epsilon, 
        SUM(con_delta) AS con_delta
        FROM (
            SELECT dataset, 
            COALESCE(all_epsilon, 0) AS all_epsilon, 
            COALESCE(all_delta, 0) as all_delta, 
            COALESCE(con_epsilon, 0) as con_epsilon, 
            COALESCE(con_delta,  0) as con_delta FROM UserBudgetAllocation
            )
         GROUP BY (dataset)
    )
    SELECT D.id, 
    D.total_epsilon, 
    COALESCE(D.total_delta, 0) AS total_delta, 
    COALESCE(C.all_epsilon, 0) AS all_epsilon, 
    COALESCE(C.all_delta, 0) AS all_delta, 
    COALESCE(C.con_epsilon, 0) AS con_epsilon, 
    COALESCE(C.con_delta, 0) AS con_delta
    FROM Dataset as D LEFT OUTER JOIN Consumed as C ON D.id = C.dataset
);

CREATE VIEW GetUsers AS (
    WITH Uroles AS (
        SELECT username, ARRAY_AGG(role) AS roles FROM UserRoles
        GROUP BY username
    )
    SELECT handle, name, roles, created_time, updated_time 
    FROM Users LEFT JOIN Uroles ON handle = username
);

CREATE VIEW GetUserBudgets AS (
    SELECT D.id,
    userid,
    COALESCE(all_epsilon, 0) AS all_epsilon,
    COALESCE(all_delta, 0) AS all_delta,
    COALESCE(con_epsilon, 0) AS con_epsilon,
    COALESCE(con_delta, 0) AS con_delta
    FROM Dataset AS D LEFT OUTER JOIN UserBudgetAllocation
    ON D.id = dataset
);
