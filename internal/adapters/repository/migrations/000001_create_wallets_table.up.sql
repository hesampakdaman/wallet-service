CREATE TABLE IF NOT EXISTS wallets (
   id UUID PRIMARY KEY,
   player_id UUID UNIQUE NOT NULL,
   balance INTEGER NOT NULL CONSTRAINT wallets_balance_non_negative CHECK (balance >= 0)
);
