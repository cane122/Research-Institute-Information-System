ALTER TABLE zahtevipromenefaze
    ALTER COLUMN zadatak_id DROP NOT NULL;
ALTER TABLE zahtevipromenefaze
    ADD COLUMN dokument_id INT NULL;
ALTER TABLE zahtevipromenefaze
    ADD CONSTRAINT fk_zahtevi_dokument_id FOREIGN KEY (dokument_id) REFERENCES dokumenti(dokument_id) ON DELETE CASCADE;