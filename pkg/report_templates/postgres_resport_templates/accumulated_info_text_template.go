package postgres_resport_templates

const AccumulatedInfoTemplate = `
# Current date: {{.CurrentDate}}
# Hostname: {{.Hostname}}
# Files:
{{.Files}}
# Overall: {{.TotalQueryCount}}, Unique: {{.TotalUniqueQueryCount}} , QPS: {{.TotalQPS}}
# Time range: {{.FromTimestamp}} to {{.ToTimestamp}}
# Attribute          total     min     max     avg     95%  stddev  median
# ============     ======= ======= ======= ======= ======= ======= =======
# Exec time        {{.Duration.Total}} {{.Duration.Min}} {{.Duration.Max}} {{.Duration.Avg}} {{.Duration.Percentile95}} {{.Duration.StdDev}} {{.Duration.Median}}
`
