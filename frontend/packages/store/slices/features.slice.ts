import { createSlice, PayloadAction } from '@reduxjs/toolkit'

export interface Features {
    chat?: boolean
    hrm?: boolean
    tms?: boolean
    meetings?: boolean
}

const initialState: Features = {}

const featureFlagsSlice = createSlice({
    name: 'features',
    initialState,
    reducers: {
        setFeatures(_, action: PayloadAction<Features>) {
            return action.payload
        },
    },
})

export const { setFeatures } = featureFlagsSlice.actions
export default featureFlagsSlice.reducer
