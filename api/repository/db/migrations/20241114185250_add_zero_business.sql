-- +goose Up
-- +goose StatementBegin
INSERT INTO businesses (id, name, created_at)
VALUES (0, 'clients', now());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE
FROM businesses
WHERE id = 0;
-- +goose StatementEnd
