import {create} from 'zustand'
import {User} from '../lib/types'

interface UserInfoState {
  user: User | null
  setUser: (user: User) => void
  clearUser: () => void
}

export const useUserInfo = create<UserInfoState>(set => ({
  user: null,
  setUser: (user: User) => set({user: user}),
  clearUser: () => set({user: null}),
}))
