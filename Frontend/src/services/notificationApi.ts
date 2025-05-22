import axios from "../lib/axios";
import { SenderEmail } from "../lib/types";


export async function getSenderEmail() {
    const response = await axios.get('/api/notification-config/sender')

    if (!response) throw new Error('Something went wrong!')

    return response.data

}

export async function setSenderEmail(data: SenderEmail) {
    const response = await axios.post('/api/notification-config/sender', data)

    if (!response) throw new Error('Something went wrong!')
    
    return response.status
}

export async function getRecipientEmail() {
    const response = await axios.get('/api/notification-config/')

    if (!response) throw new Error('Something went wrong!')

    return response.data
}