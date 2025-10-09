-- Analytics Dummy Data
-- Populate LogAktivnosti table with realistic test data

-- Clear existing test data (optional)
-- DELETE FROM LogAktivnosti;

-- Insert dummy activities for the last 30 days
DO $$
DECLARE
    doc_id INT;
    user_id INT;
    activity_date TIMESTAMP;
    activity_count INT := 0;
BEGIN
    -- Generate activities for each day in the last 30 days
    FOR day_offset IN 0..29 LOOP
        activity_date := CURRENT_TIMESTAMP - (day_offset || ' days')::INTERVAL;
        
        -- Generate 5-15 random activities per day
        FOR i IN 1..(5 + floor(random() * 10))::INT LOOP
            -- Random user (1-5)
            user_id := 1 + floor(random() * 5)::INT;
            
            -- Random document (1-20)
            doc_id := 1 + floor(random() * 20)::INT;
            
            -- Random activity type
            CASE floor(random() * 6)::INT
                WHEN 0 THEN
                    -- UPLOAD activity
                    INSERT INTO LogAktivnosti (
                        korisnik_id, tip_aktivnosti, ciljani_entitet, ciljani_id, 
                        opis, datuma
                    ) VALUES (
                        user_id, 'UPLOAD', 'DOKUMENT', doc_id,
                        'Uploaded document: Document_' || doc_id || '.pdf',
                        activity_date - (random() * INTERVAL '24 hours')
                    );
                    
                WHEN 1 THEN
                    -- VIEW activity (most common)
                    INSERT INTO LogAktivnosti (
                        korisnik_id, tip_aktivnosti, ciljani_entitet, ciljani_id, 
                        opis, datuma
                    ) VALUES (
                        user_id, 'VIEW', 'DOKUMENT', doc_id,
                        'Viewed document: Document_' || doc_id || '.pdf',
                        activity_date - (random() * INTERVAL '24 hours')
                    );
                    
                WHEN 2 THEN
                    -- EDIT activity
                    INSERT INTO LogAktivnosti (
                        korisnik_id, tip_aktivnosti, ciljani_entitet, ciljani_id, 
                        opis, datuma
                    ) VALUES (
                        user_id, 'EDIT', 'DOKUMENT', doc_id,
                        'Updated document: Document_' || doc_id || '.pdf',
                        activity_date - (random() * INTERVAL '24 hours')
                    );
                    
                WHEN 3 THEN
                    -- DELETE activity
                    INSERT INTO LogAktivnosti (
                        korisnik_id, tip_aktivnosti, ciljani_entitet, ciljani_id, 
                        opis, datuma
                    ) VALUES (
                        user_id, 'DELETE', 'DOKUMENT', doc_id,
                        'Deleted document: Document_' || doc_id || '.pdf',
                        activity_date - (random() * INTERVAL '24 hours')
                    );
                    
                WHEN 4 THEN
                    -- LOGIN activity
                    INSERT INTO LogAktivnosti (
                        korisnik_id, tip_aktivnosti, ciljani_entitet, ciljani_id, 
                        opis, datuma
                    ) VALUES (
                        user_id, 'LOGIN', 'KORISNIK', user_id,
                        'User logged in: user' || user_id,
                        activity_date - (random() * INTERVAL '24 hours')
                    );
                    
                WHEN 5 THEN
                    -- LOGOUT activity
                    INSERT INTO LogAktivnosti (
                        korisnik_id, tip_aktivnosti, ciljani_entitet, ciljani_id, 
                        opis, datuma
                    ) VALUES (
                        user_id, 'LOGOUT', 'KORISNIK', user_id,
                        'User logged out: user' || user_id,
                        activity_date - (random() * INTERVAL '24 hours')
                    );
            END CASE;
            
            activity_count := activity_count + 1;
        END LOOP;
    END LOOP;
    
    RAISE NOTICE 'Inserted % activity records', activity_count;
END $$;

-- Add some extra VIEW activities to make the count higher
INSERT INTO LogAktivnosti (korisnik_id, tip_aktivnosti, ciljani_entitet, ciljani_id, opis, datuma)
SELECT 
    1 + floor(random() * 5)::INT,
    'VIEW',
    'DOKUMENT',
    1 + floor(random() * 20)::INT,
    'Viewed document: Document_' || (1 + floor(random() * 20)::INT) || '.pdf',
    CURRENT_TIMESTAMP - (random() * INTERVAL '30 days')
FROM generate_series(1, 2000); -- Generate 2000 additional VIEW activities

-- Verify the data
SELECT 
    tip_aktivnosti,
    COUNT(*) as broj_aktivnosti,
    MIN(datuma) as najstarija,
    MAX(datuma) as najnovija
FROM LogAktivnosti
GROUP BY tip_aktivnosti
ORDER BY broj_aktivnosti DESC;

-- Show summary
SELECT 
    'Total Activities' as metric,
    COUNT(*)::TEXT as value
FROM LogAktivnosti
UNION ALL
SELECT 
    'Total Users' as metric,
    COUNT(DISTINCT korisnik_id)::TEXT as value
FROM LogAktivnosti
UNION ALL
SELECT 
    'Total Documents' as metric,
    COUNT(DISTINCT ciljani_id)::TEXT as value
FROM LogAktivnosti
WHERE ciljani_entitet = 'DOKUMENT';
