import { Message } from '@/packages/types'
import { createSlice, PayloadAction } from '@reduxjs/toolkit'

export interface ChatState {
    messages: Message[]
}

const initialState: ChatState = { messages: [] }

const chatSlice = createSlice({
    name: 'chat',
    initialState,
    reducers: {
        addMessage(state, action: PayloadAction<Message>) {
            state.messages.push(action.payload)
        },
        setMessages(state, action: PayloadAction<Message[]>) {
            state.messages = action.payload
        },
    },
})

export const { addMessage, setMessages } = chatSlice.actions
export default chatSlice.reducer
