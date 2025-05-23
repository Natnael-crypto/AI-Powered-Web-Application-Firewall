import { useMutation, useQuery } from "@tanstack/react-query"
import { getRecipientEmail, getSenderEmail, setSenderEmail } from "../../services/notificationApi"

export function useGetSenderEmail() {
    return useQuery({
      queryKey: ['senderEmail'],
      queryFn: getSenderEmail,
    })
  }

export function useSetSenderEmail() {
    return useMutation({
        mutationKey: ['setSenderEmail'],
        mutationFn: setSenderEmail,
    })
}

export function useGetRecipientEmail() {
    return useQuery({
        queryKey: ['recipientEmail'],
        queryFn: getRecipientEmail,
    })
}