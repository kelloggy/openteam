// tasks/04‑sql‑reasoning/go/queries.go
package queries

// Task A
const SQLA = `SELECT
    c.id AS campaign_id,
    COALESCE(SUM(p.amount_thb), 0) AS total_thb,
    ROUND(COALESCE(SUM(p.amount_thb), 0) * 1.0 / c.target_thb, 4) AS pct_of_target
FROM
    campaign c
LEFT JOIN
    pledge p ON c.id = p.campaign_id
GROUP BY
    c.id, c.target_thb
ORDER BY
    pct_of_target DESC,
    campaign_id ASC;`

// Task B
const SQLB = `
WITH
thailand AS (
  SELECT p.amount_thb
  FROM pledge p
  JOIN donor d ON p.donor_id = d.id
  WHERE d.country = 'Thailand'
  ORDER BY p.amount_thb
),


thailand_ranked AS (
  SELECT
    amount_thb,
    ROW_NUMBER() OVER () AS rn,
    (SELECT COUNT(*) FROM thailand) AS total
  FROM thailand
),

global AS (
  SELECT amount_thb
  FROM pledge
  ORDER BY amount_thb
),

global_ranked AS (
  SELECT
    amount_thb,
    ROW_NUMBER() OVER () AS rn,
    (SELECT COUNT(*) FROM global) AS total
  FROM global
)

SELECT 'global' AS scope, amount_thb AS p90_thb
FROM global_ranked
WHERE rn = CAST(0.9 * total + 0.9999 AS INT)

UNION ALL

SELECT 'thailand' AS scope, amount_thb AS p90_thb
FROM thailand_ranked
WHERE rn = CAST(0.9 * total + 0.9999 AS INT)

ORDER BY scope;
`

var Indexes = []string{
	"donor.country",
	"pledge.amount_thb",
} // skipped
