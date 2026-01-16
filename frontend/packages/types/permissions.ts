export interface Permission {
    action: string
    resource: string
    conditions?: Record<string, any>
}
