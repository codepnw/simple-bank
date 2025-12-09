-- Clear data reset id = 1
TRUNCATE TABLE users, accounts, entries, transfers RESTART IDENTITY CASCADE;

INSERT INTO users (email, username, password, first_name, last_name) VALUES
('user1@example.com', 'somchai', '$2a$10$2Cgf4hs0BxKCbhQwt2Tq/euFD9FWd0WMicPEs/7nukVZzIniE3.Om', 'Somchai', 'Jaidee'),
('user2@example.com', 'somsri', '$2a$10$2Cgf4hs0BxKCbhQwt2Tq/euFD9FWd0WMicPEs/7nukVZzIniE3.Om', 'Somsri', 'Rakdee'),
('rich@example.com', 'elon', '$2a$10$2Cgf4hs0BxKCbhQwt2Tq/euFD9FWd0WMicPEs/7nukVZzIniE3.Om', 'Elon', 'Musk');

INSERT INTO accounts (owner_id, balance, currency) VALUES
(1, 500000, 'THB'), -- 5,000.00
(1, 10000, 'USD'), -- 100.00
(2, 0, 'THB'), -- 0.00
(3, 100000000, 'THB'); -- 1,000,000.00
