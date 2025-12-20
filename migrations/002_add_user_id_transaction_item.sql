ALTER TABLE transaction_items 
ADD COLUMN user_id VARCHAR(36);

UPDATE transaction_items ti
SET user_id = t.user_id
FROM transactions t
WHERE ti.transaction_id = t.id;

ALTER TABLE transaction_items 
ALTER COLUMN user_id SET NOT NULL;

ALTER TABLE transaction_items
ADD CONSTRAINT fk_transaction_items_user
FOREIGN KEY (user_id) REFERENCES users(id);

CREATE INDEX idx_transaction_items_user_id ON transaction_items(user_id);