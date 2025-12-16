-- ============================================================================
-- Populate Database with Realistic Data
-- ============================================================================

-- Add more projects
INSERT INTO projekti (naziv_projekta, opis, datum_pocetka, datum_zavrsetka, status, rukovodilac_id)
VALUES ('AI Research Initiative', 'Development of machine learning algorithms for data analysis', 
        TO_DATE('2024-01-15', 'YYYY-MM-DD'), TO_DATE('2025-12-31', 'YYYY-MM-DD'), 'U toku', 43);

INSERT INTO projekti (naziv_projekta, opis, datum_pocetka, datum_zavrsetka, status, rukovodilac_id)
VALUES ('Cloud Infrastructure Migration', 'Migration of on-premise systems to cloud infrastructure', 
        TO_DATE('2024-03-01', 'YYYY-MM-DD'), TO_DATE('2025-06-30', 'YYYY-MM-DD'), 'U toku', 44);

INSERT INTO projekti (naziv_projekta, opis, datum_pocetka, datum_zavrsetka, status, rukovodilac_id)
VALUES ('Cybersecurity Assessment', 'Comprehensive security audit and vulnerability assessment', 
        TO_DATE('2024-02-10', 'YYYY-MM-DD'), TO_DATE('2024-12-15', 'YYYY-MM-DD'), 'Završen', 43);

INSERT INTO projekti (naziv_projekta, opis, datum_pocetka, datum_zavrsetka, status, rukovodilac_id)
VALUES ('Data Analytics Platform', 'Building enterprise-wide analytics and reporting platform', 
        TO_DATE('2024-04-01', 'YYYY-MM-DD'), TO_DATE('2026-03-31', 'YYYY-MM-DD'), 'U toku', 44);

INSERT INTO projekti (naziv_projekta, opis, datum_pocetka, datum_zavrsetka, status, rukovodilac_id)
VALUES ('Mobile App Development', 'Cross-platform mobile application for research collaboration', 
        TO_DATE('2024-05-15', 'YYYY-MM-DD'), TO_DATE('2025-08-15', 'YYYY-MM-DD'), 'U toku', 2);

-- Add more documents
INSERT INTO dokumenti (projekat_id, naziv_dokumenta, opis, tip_dokumenta, jezik_dokumenta, kreirao_korisnik_id)
VALUES (1, 'Project Charter - AI Initiative', 'Initial project planning document', 'PDF', 'English', 43);

INSERT INTO dokumenti (projekat_id, naziv_dokumenta, opis, tip_dokumenta, jezik_dokumenta, kreirao_korisnik_id)
VALUES (1, 'Technical Specification v1.0', 'Detailed technical requirements and architecture', 'DOCX', 'English', 42);

INSERT INTO dokumenti (projekat_id, naziv_dokumenta, opis, tip_dokumenta, jezik_dokumenta, kreirao_korisnik_id)
VALUES (2, 'Cloud Migration Strategy', 'Strategy document for cloud migration', 'PDF', 'English', 44);

INSERT INTO dokumenti (projekat_id, naziv_dokumenta, opis, tip_dokumenta, jezik_dokumenta, kreirao_korisnik_id)
VALUES (2, 'Risk Assessment Report', 'Identified risks and mitigation strategies', 'DOCX', 'English', 45);

INSERT INTO dokumenti (projekat_id, naziv_dokumenta, opis, tip_dokumenta, jezik_dokumenta, kreirao_korisnik_id)
VALUES (3, 'Security Audit Report', 'Comprehensive security assessment findings', 'PDF', 'English', 43);

INSERT INTO dokumenti (projekat_id, naziv_dokumenta, opis, tip_dokumenta, jezik_dokumenta, kreirao_korisnik_id)
VALUES (4, 'Analytics Platform Design', 'System architecture and design document', 'PDF', 'English', 44);

INSERT INTO dokumenti (projekat_id, naziv_dokumenta, opis, tip_dokumenta, jezik_dokumenta, kreirao_korisnik_id)
VALUES (5, 'Mobile App Wireframes', 'UI/UX design and wireframes', 'PDF', 'English', 2);

INSERT INTO dokumenti (projekat_id, naziv_dokumenta, opis, tip_dokumenta, jezik_dokumenta, kreirao_korisnik_id)
VALUES (1, 'Meeting Minutes 2024-Q1', 'Quarterly project review meeting notes', 'DOCX', 'Serbian', 42);

INSERT INTO dokumenti (projekat_id, naziv_dokumenta, opis, tip_dokumenta, jezik_dokumenta, kreirao_korisnik_id)
VALUES (2, 'Budget Proposal', 'Detailed budget breakdown and justification', 'XLSX', 'English', 44);

INSERT INTO dokumenti (projekat_id, naziv_dokumenta, opis, tip_dokumenta, jezik_dokumenta, kreirao_korisnik_id)
VALUES (3, 'Compliance Checklist', 'Security compliance requirements checklist', 'PDF', 'English', 46);

-- Grant permissions to manager2 (user 44) for viewing and editing documents
-- Manager2 can view most documents
INSERT INTO dozvoledokumenata (dokument_id, korisnik_id, moze_citati, moze_menjati, moze_brisati)
VALUES (3, 44, 1, 0, 0);

INSERT INTO dozvoledokumenata (dokument_id, korisnik_id, moze_citati, moze_menjati, moze_brisati)
VALUES (21, 44, 1, 1, 0);

INSERT INTO dozvoledokumenata (dokument_id, korisnik_id, moze_citati, moze_menjati, moze_brisati)
VALUES (22, 44, 1, 1, 0);

INSERT INTO dozvoledokumenata (dokument_id, korisnik_id, moze_citati, moze_menjati, moze_brisati)
VALUES (50, 44, 1, 0, 0);

INSERT INTO dozvoledokumenata (dokument_id, korisnik_id, moze_citati, moze_menjati, moze_brisati)
VALUES (51, 44, 1, 1, 0);

INSERT INTO dozvoledokumenata (dokument_id, korisnik_id, moze_citati, moze_menjati, moze_brisati)
VALUES (52, 44, 1, 1, 0);

-- Add realistic activity logs
-- Login activities
INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (1, 'LOGIN', 'KORISNIK', 1, 'Administrator logged in');

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (44, 'LOGIN', 'KORISNIK', 44, 'manager2 logged in');

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (43, 'LOGIN', 'KORISNIK', 43, 'manager1 logged in');

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (42, 'LOGIN', 'KORISNIK', 42, 'researcher1 logged in');

-- Document view activities
INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (44, 'VIEW', 'DOKUMENT', 3, 'Viewed document: Test_Dokument.pdf');

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (44, 'VIEW', 'DOKUMENT', 21, 'Viewed document: Studija_A.pdf');

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (43, 'VIEW', 'DOKUMENT', 50, 'Viewed document: Uputstvo za koriscenje sistema');

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (42, 'VIEW', 'DOKUMENT', 51, 'Viewed document: Pravilnik o zastiti podataka');

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (1, 'VIEW', 'DOKUMENT', 52, 'Viewed document: Sablon istrazivackog izvestaja');

-- Document edit activities
INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (44, 'EDIT', 'DOKUMENT', 21, 'Updated document: Studija_A.pdf');

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (44, 'EDIT', 'DOKUMENT', 22, 'Updated document: Plan_B.docx');

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (1, 'EDIT', 'DOKUMENT', 50, 'Updated document metadata');

-- Document upload activities
INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (43, 'UPLOAD', 'DOKUMENT', 60, 'Uploaded new document: Project Charter - AI Initiative');

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (42, 'UPLOAD', 'DOKUMENT', 61, 'Uploaded new document: Technical Specification v1.0');

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (44, 'UPLOAD', 'DOKUMENT', 62, 'Uploaded new document: Cloud Migration Strategy');

-- Project creation activities
INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (43, 'CREATE', 'PROJEKAT', 1, 'Created new project: AI Research Initiative');

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (44, 'CREATE', 'PROJEKAT', 2, 'Created new project: Cloud Infrastructure Migration');

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (43, 'CREATE', 'PROJEKAT', 3, 'Created new project: Cybersecurity Assessment');

-- Project view activities
INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (44, 'VIEW', 'PROJEKAT', 1, 'Viewed project details: AI Research Initiative');

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (43, 'VIEW', 'PROJEKAT', 2, 'Viewed project details: Cloud Infrastructure Migration');

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (42, 'VIEW', 'PROJEKAT', 3, 'Viewed project details: Cybersecurity Assessment');

-- Project update activities
INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (43, 'EDIT', 'PROJEKAT', 1, 'Updated project status');

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (44, 'EDIT', 'PROJEKAT', 2, 'Updated project timeline');

-- Permission grant activities
INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (1, 'PERMISSION_GRANT', 'DOKUMENT', 21, 'Granted edit permission to manager2');

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (1, 'PERMISSION_GRANT', 'DOKUMENT', 22, 'Granted edit permission to manager2');

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (43, 'PERMISSION_GRANT', 'DOKUMENT', 51, 'Granted view permission to manager2');

-- Search activities
INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (44, 'SEARCH', NULL, NULL, 'Searched for documents containing: security');

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (42, 'SEARCH', NULL, NULL, 'Searched for documents containing: research');

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (43, 'SEARCH', NULL, NULL, 'Searched for projects containing: AI');

-- Export activities
INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (44, 'EXPORT', 'DOKUMENT', 21, 'Exported document: Studija_A.pdf');

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (43, 'EXPORT', 'PROJEKAT', 1, 'Exported project report');

-- Logout activities
INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (44, 'LOGOUT', 'KORISNIK', 44, 'manager2 logged out');

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (43, 'LOGOUT', 'KORISNIK', 43, 'manager1 logged out');

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis)
VALUES (42, 'LOGOUT', 'KORISNIK', 42, 'researcher1 logged out');

-- Add tasks for projects
INSERT INTO zadaci (projekat_id, naziv_zadatka, opis, prioritet, status, dodeljen_korisniku_id, rok_za_zavrsetak)
VALUES (1, 'Data Collection and Preprocessing', 'Gather and prepare training data for ML models', 'Visok', 'U toku', 42, TO_DATE('2025-02-15', 'YYYY-MM-DD'));

INSERT INTO zadaci (projekat_id, naziv_zadatka, opis, prioritet, status, dodeljen_korisniku_id, rok_za_zavrsetak)
VALUES (1, 'Model Training and Validation', 'Train and validate machine learning models', 'Visok', 'Zakazan', 45, TO_DATE('2025-03-30', 'YYYY-MM-DD'));

INSERT INTO zadaci (projekat_id, naziv_zadatka, opis, prioritet, status, dodeljen_korisniku_id, rok_za_zavrsetak)
VALUES (2, 'Infrastructure Assessment', 'Evaluate current infrastructure and requirements', 'Srednji', 'Završen', 46, TO_DATE('2024-04-15', 'YYYY-MM-DD'));

INSERT INTO zadaci (projekat_id, naziv_zadatka, opis, prioritet, status, dodeljen_korisniku_id, rok_za_zavrsetak)
VALUES (2, 'Cloud Provider Selection', 'Compare and select optimal cloud provider', 'Visok', 'U toku', 44, TO_DATE('2025-01-31', 'YYYY-MM-DD'));

INSERT INTO zadaci (projekat_id, naziv_zadatka, opis, prioritet, status, dodeljen_korisniku_id, rok_za_zavrsetak)
VALUES (3, 'Vulnerability Scanning', 'Perform comprehensive security scans', 'Kritičan', 'Završen', 47, TO_DATE('2024-11-30', 'YYYY-MM-DD'));

INSERT INTO zadaci (projekat_id, naziv_zadatka, opis, prioritet, status, dodeljen_korisniku_id, rok_za_zavrsetak)
VALUES (4, 'Database Design', 'Design analytics data warehouse schema', 'Visok', 'U toku', 48, TO_DATE('2025-02-28', 'YYYY-MM-DD'));

INSERT INTO zadaci (projekat_id, naziv_zadatka, opis, prioritet, status, dodeljen_korisniku_id, rok_za_zavrsetak)
VALUES (5, 'Mobile UI Design', 'Create mobile application user interface', 'Srednji', 'U toku', 49, TO_DATE('2025-01-20', 'YYYY-MM-DD'));

-- Add more recent activity logs (last 30 days)
INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis, kreiran_datuma)
VALUES (44, 'LOGIN', 'KORISNIK', 44, 'manager2 logged in', SYSTIMESTAMP - INTERVAL '1' DAY);

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis, kreiran_datuma)
VALUES (44, 'VIEW', 'DOKUMENT', 50, 'Viewed system documentation', SYSTIMESTAMP - INTERVAL '1' DAY);

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis, kreiran_datuma)
VALUES (44, 'EDIT', 'DOKUMENT', 21, 'Updated project documentation', SYSTIMESTAMP - INTERVAL '1' DAY);

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis, kreiran_datuma)
VALUES (44, 'VIEW', 'PROJEKAT', 2, 'Reviewed cloud migration progress', SYSTIMESTAMP - INTERVAL '1' DAY);

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis, kreiran_datuma)
VALUES (1, 'LOGIN', 'KORISNIK', 1, 'Administrator logged in', SYSTIMESTAMP - INTERVAL '2' HOUR);

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis, kreiran_datuma)
VALUES (1, 'VIEW', 'DOKUMENT', 55, 'Viewed system guidelines', SYSTIMESTAMP - INTERVAL '1' HOUR);

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis, kreiran_datuma)
VALUES (42, 'LOGIN', 'KORISNIK', 42, 'researcher1 logged in', SYSTIMESTAMP - INTERVAL '3' HOUR);

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis, kreiran_datuma)
VALUES (42, 'VIEW', 'PROJEKAT', 1, 'Checked AI project tasks', SYSTIMESTAMP - INTERVAL '2' HOUR);

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis, kreiran_datuma)
VALUES (43, 'LOGIN', 'KORISNIK', 43, 'manager1 logged in', SYSTIMESTAMP - INTERVAL '5' HOUR);

INSERT INTO logaktivnosti (korisnik_id, tip_aktivnosti, entitet_tip, entitet_id, opis, kreiran_datuma)
VALUES (43, 'EDIT', 'PROJEKAT', 1, 'Updated AI project timeline', SYSTIMESTAMP - INTERVAL '4' HOUR);

COMMIT;

-- Display summary
SELECT 'Added projects:' as summary, COUNT(*) as count FROM projekti WHERE projekat_id > 5
UNION ALL
SELECT 'Added documents:', COUNT(*) FROM dokumenti WHERE dokument_id > 59
UNION ALL
SELECT 'Added permissions for manager2:', COUNT(*) FROM dozvoledokumenata WHERE korisnik_id = 44
UNION ALL
SELECT 'Added activity logs:', COUNT(*) FROM logaktivnosti WHERE log_id > 20
UNION ALL
SELECT 'Added tasks:', COUNT(*) FROM zadaci WHERE zadatak_id > 5;
