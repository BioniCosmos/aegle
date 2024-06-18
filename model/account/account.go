package account

type Role string

const Member, Admin Role = "member", "admin"

type Status string

const Normal, Unverified Status = "normal", "unverified"
