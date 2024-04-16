CREATE TABLE IF NOT EXISTS customers (
    id varchar(255),
    document_id varchar(255),
    document_type int,
    is_anonymous boolean,
    password varchar(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id)
);