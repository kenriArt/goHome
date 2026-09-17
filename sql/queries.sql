INSERT INTO orders (id, customer, address)
VALUES (1, 'Vasya', 'ул. Лесная, 12')

SELECT * FROM orders

UPDATE orders
SET is_delivered = true
WHERE id = 1

DELETE FROM orders
WHERE id = 1

SELECT COUNT(*)
FROM orders
WHERE is_delivered = true
