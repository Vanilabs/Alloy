import { combineReducers } from "@reduxjs/toolkit";
import { auth, chat, features, runtime, tenant } from "../slices";

export const rootReducer = combineReducers({
    runtime,
    auth,
    chat,
    features,
    tenant,
});
