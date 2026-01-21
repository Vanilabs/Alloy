import { AlloyModule } from '@/packages/shared/types'
import { createSlice, PayloadAction } from '@reduxjs/toolkit'


const initialState: AlloyModule[] = []

export const runtimeSlice = createSlice({
    name: 'runtime',
    initialState,
    reducers: {
        registerModule(state, action: PayloadAction<AlloyModule>) {
            state.push(action.payload)
        },
        moduleRegistered(state, action: PayloadAction<AlloyModule>) {
            if (!state.includes(action.payload)) {
                state.push(action.payload)
            }
        },
        clearModules(state) {
            return initialState
        }
    },
    extraReducers: builder => {},
})

export const { moduleRegistered, clearModules } = runtimeSlice.actions
export default runtimeSlice.reducer