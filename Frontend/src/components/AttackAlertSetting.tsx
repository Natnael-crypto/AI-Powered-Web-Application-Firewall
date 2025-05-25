import {useEffect, useState} from 'react'
import {Info} from 'lucide-react'
import {useForm} from 'react-hook-form'
import {zodResolver} from '@hookform/resolvers/zod'
import {SenderEmail, SenderEmailSchema} from '../lib/types'
import {useGetSenderEmail, useSetSenderEmail} from '../hooks/api/useNotification'
import {useToast} from '../hooks/useToast'

const AttackAlertSettings = () => {
  const [alertType, setAlertType] = useState<'Telegram' | 'Email'>('Email')
  // const [webhook, setWebhook] = useState('')

  const {mutate: setSenderEmail, data: statusCode} = useSetSenderEmail()
  const {data: senderEmailConfig, refetch: refetchSenderEmail} = useGetSenderEmail()
  const {addToast: toast} = useToast()

  const {
    register,
    handleSubmit,
    setValue,
    formState: {errors},
  } = useForm<SenderEmail>({
    resolver: zodResolver(SenderEmailSchema), // Use Zod resolver for validation
  })

  const onSave = (data: SenderEmail) => {
    setSenderEmail(data, {
      onSuccess: () => {
        toast('Email set successfully')
        refetchSenderEmail()
      },
      onError: () => {
        toast('Something went wrong while setting up sender email')
      },
    })
  }

  useEffect(() => {
    if (senderEmailConfig) {
      setValue('sender_email', senderEmailConfig.sender_email)
      setValue('app_password', senderEmailConfig.app_password)
    }
  }, [senderEmailConfig, setValue])

  return (
    <div className="p-6 bg-white  shadow-lg w-full">
      {/* Title */}
      <div className="flex items-center mb-4">
        <h4 className="text-lg font-semibold text-gray-800 mr-2">Attack Alert</h4>
        <Info size={16} className="text-blue-500" />
      </div>

      {/* Radio buttons */}
      <div className="flex gap-6 mb-4">
        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="radio"
            name="alertType"
            value="Email"
            checked={alertType === 'Email'}
            onChange={() => setAlertType('Email')}
            className="w-4 h-4 text-blue-600 border-gray-300 focus:ring-blue-500"
          />
          <span className="text-gray-700">Email</span>
        </label>
      </div>

      {/* Webhook input
      <input
        type="text"
        placeholder="https://api.xxx.com/robot/send?access_token=xxxxxx"
        value={webhook}
        onChange={e => setWebhook(e.target.value)}
        className="w-full px-4 py-2 mb-4 border border-gray-300  text-sm placeholder-gray-400"
      /> */}

      <form onSubmit={handleSubmit(onSave)} className="space-y-6">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Sender Email
            </label>
            <input
              type="email"
              {...register('sender_email')}
              placeholder="Enter sender email"
              className="w-full p-3 border border-gray-300 shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            {errors.sender_email && (
              <p className="text-red-500">{errors.sender_email.message}</p>
            )}
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              App Password
            </label>
            <input
              type="password"
              {...register('app_password')}
              placeholder="Enter your Gmail or Outlook app password"
              className="w-full p-3 border border-gray-300 shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            {errors.app_password && (
              <p className="text-red-500">{errors.app_password.message}</p>
            )}
          </div>
        </div>
        <div className="mt-6 flex justify-end">
          <button
            type="submit"
            className="bg-black text-white font-semibold px-6 py-2 hover:bg-slate-800 transition"
          >
            Save
          </button>
        </div>
      </form>
    </div>
  )
}

export default AttackAlertSettings
