CREATE TABLE IF NOT EXISTS chat_messages (
    room_id uuid,
    chat_message_id timeuuid,
    author_id uuid,
    content text,

    PRIMARY KEY (room_id, chat_message_id)
);
