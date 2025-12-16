-- Update existing documents 70-75 to assign to manager2's projects
-- manager2 manages: Cloud Infrastructure Migration (49), Data Analytics Platform (51)

-- Update documents for Cloud Infrastructure Migration project (49)
UPDATE dokumenti SET projekat_id = 49, naziv_dokumenta = 'Cloud Migration Strategy', 
       opis = 'Comprehensive strategy for migrating infrastructure to cloud', 
       tip_dokumenta = 'Plan', jezik_dokumenta = 'EN'
WHERE dokument_id = 70;

UPDATE dokumenti SET projekat_id = 49, naziv_dokumenta = 'Infrastructure Assessment Report', 
       opis = 'Current infrastructure audit and requirements analysis', 
       tip_dokumenta = 'Report', jezik_dokumenta = 'EN'
WHERE dokument_id = 71;

UPDATE dokumenti SET projekat_id = 49, naziv_dokumenta = 'Cloud Cost Estimation', 
       opis = 'Projected costs for cloud infrastructure', 
       tip_dokumenta = 'Spreadsheet', jezik_dokumenta = 'EN'
WHERE dokument_id = 72;

-- Update documents for Data Analytics Platform project (51)
UPDATE dokumenti SET projekat_id = 51, naziv_dokumenta = 'Data Analytics Architecture', 
       opis = 'System architecture for analytics platform', 
       tip_dokumenta = 'Technical', jezik_dokumenta = 'EN'
WHERE dokument_id = 73;

UPDATE dokumenti SET projekat_id = 51, naziv_dokumenta = 'Analytics Data Model', 
       opis = 'Data warehouse schema and ETL processes', 
       tip_dokumenta = 'Technical', jezik_dokumenta = 'EN'
WHERE dokument_id = 74;

UPDATE dokumenti SET projekat_id = 51, naziv_dokumenta = 'BI Dashboard Requirements', 
       opis = 'Business Intelligence dashboard specifications', 
       tip_dokumenta = 'Requirements', jezik_dokumenta = 'EN'
WHERE dokument_id = 75;

-- Add tags for documents
DELETE FROM dokumenttagovi WHERE dokument_id BETWEEN 70 AND 75;

INSERT INTO tagovi (tag_id, naziv_taga) 
SELECT 100, 'cloud' FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM tagovi WHERE naziv_taga = 'cloud');

INSERT INTO tagovi (tag_id, naziv_taga) 
SELECT 101, 'migration' FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM tagovi WHERE naziv_taga = 'migration');

INSERT INTO tagovi (tag_id, naziv_taga) 
SELECT 102, 'infrastructure' FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM tagovi WHERE naziv_taga = 'infrastructure');

INSERT INTO tagovi (tag_id, naziv_taga) 
SELECT 103, 'analytics' FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM tagovi WHERE naziv_taga = 'analytics');

INSERT INTO tagovi (tag_id, naziv_taga) 
SELECT 104, 'cost' FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM tagovi WHERE naziv_taga = 'cost');

-- Link tags to documents
INSERT INTO dokumenttagovi (dokument_id, tag_id) 
SELECT 70, tag_id FROM tagovi WHERE naziv_taga IN ('cloud', 'migration');

INSERT INTO dokumenttagovi (dokument_id, tag_id) 
SELECT 71, tag_id FROM tagovi WHERE naziv_taga IN ('infrastructure', 'migration');

INSERT INTO dokumenttagovi (dokument_id, tag_id) 
SELECT 72, tag_id FROM tagovi WHERE naziv_taga IN ('cloud', 'cost');

INSERT INTO dokumenttagovi (dokument_id, tag_id) 
SELECT 73, tag_id FROM tagovi WHERE naziv_taga IN ('analytics', 'infrastructure');

INSERT INTO dokumenttagovi (dokument_id, tag_id) 
SELECT 74, tag_id FROM tagovi WHERE naziv_taga = 'analytics';

INSERT INTO dokumenttagovi (dokument_id, tag_id) 
SELECT 75, tag_id FROM tagovi WHERE naziv_taga = 'analytics';

COMMIT;
