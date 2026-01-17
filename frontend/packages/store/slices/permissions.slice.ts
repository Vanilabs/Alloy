import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import type { Permission } from '../../types'

const initialState: Permission[] = []

const permissionsSlice = createSlice({
    name: 'permissions',
    initialState,
    reducers: {
        setPermissions(_, action: PayloadAction<Permission[]>) {
            return action.payload
        },
        clearPermissions() {
            return []
        },
    },
})

export const { setPermissions, clearPermissions } = permissionsSlice.actions
export default permissionsSlice.reducer
