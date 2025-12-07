-- Clear data reset id = 1
TRUNCATE TABLE users, accounts, entries, transfers RESTART IDENTITY CASCADE;

INSERT INTO users (email, username, password, first_name, last_name) VALUES
('user1@mail.com', 'somchai', '$2a$10$2Cgf4hs0BxKCbhQwt2Tq/euFD9FWd0WMicPEs/7nukVZzIniE3.Om', 'Somchai', 'Jaidee'),
('user2@mail.com', 'somsri', '$2a$10$2Cgf4hs0BxKCbhQwt2Tq/euFD9FWd0WMicPEs/7nukVZzIniE3.Om', 'Somsri', 'Rakdee'),
('rich@mail.com', 'elon', '$2a$10$2Cgf4hs0BxKCbhQwt2Tq/euFD9FWd0WMicPEs/7nukVZzIniE3.Om', 'Elon', 'Musk');

INSERT INTO accounts (owner_id, balance, currency) VALUES
(1, 5000, 'THB'),
(1, 100, 'USD'),
(2, 0, 'THB'),
(3, 100000, 'THB');
