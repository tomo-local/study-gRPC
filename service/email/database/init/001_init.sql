-- Email Service Database Initialization

-- Users table for authentication
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Email messages table
CREATE TABLE IF NOT EXISTS emails (
    id SERIAL PRIMARY KEY,
    message_id VARCHAR(255) UNIQUE NOT NULL,
    from_email VARCHAR(255) NOT NULL,
    to_emails TEXT[] NOT NULL,
    cc_emails TEXT[],
    bcc_emails TEXT[],
    subject VARCHAR(512),
    body_text TEXT,
    body_html TEXT,
    attachments JSONB,
    status VARCHAR(50) DEFAULT 'pending',
    sent_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Email folders/mailboxes (for IMAP)
CREATE TABLE IF NOT EXISTS mailboxes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    path VARCHAR(512) NOT NULL,
    uidvalidity INTEGER NOT NULL,
    uidnext INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, name)
);

-- Email message storage in mailboxes
CREATE TABLE IF NOT EXISTS message_mailbox (
    id SERIAL PRIMARY KEY,
    email_id INTEGER REFERENCES emails(id) ON DELETE CASCADE,
    mailbox_id INTEGER REFERENCES mailboxes(id) ON DELETE CASCADE,
    uid INTEGER NOT NULL,
    flags TEXT[] DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(mailbox_id, uid)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_emails_message_id ON emails(message_id);
CREATE INDEX IF NOT EXISTS idx_emails_from_email ON emails(from_email);
CREATE INDEX IF NOT EXISTS idx_emails_status ON emails(status);
CREATE INDEX IF NOT EXISTS idx_emails_created_at ON emails(created_at);
CREATE INDEX IF NOT EXISTS idx_mailboxes_user_id ON mailboxes(user_id);
CREATE INDEX IF NOT EXISTS idx_message_mailbox_email_id ON message_mailbox(email_id);
CREATE INDEX IF NOT EXISTS idx_message_mailbox_mailbox_id ON message_mailbox(mailbox_id);

-- Insert sample data
INSERT INTO users (email, password_hash, full_name) VALUES
('admin@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Admin User'),
('test@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Test User')
ON CONFLICT (email) DO NOTHING;

-- Create default mailboxes
INSERT INTO mailboxes (user_id, name, path, uidvalidity) VALUES
(1, 'INBOX', 'INBOX', 1),
(1, 'Sent', 'Sent', 1),
(1, 'Drafts', 'Drafts', 1),
(1, 'Trash', 'Trash', 1),
(2, 'INBOX', 'INBOX', 1),
(2, 'Sent', 'Sent', 1),
(2, 'Drafts', 'Drafts', 1),
(2, 'Trash', 'Trash', 1)
ON CONFLICT (user_id, name) DO NOTHING;
