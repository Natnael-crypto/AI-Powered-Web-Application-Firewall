import {create} from 'zustand'
import {filterOperations, logFilterType} from '../lib/types'

interface FilterState {
  filterType: logFilterType | null
  filterOperation: filterOperations | null
  selectFilterType: (type: logFilterType) => void
  selectFilterOperation: (operation: filterOperations) => void
}

export const useLogFilter = create<FilterState>(set => ({
  filterType: null,
  filterOperation: null,
  selectFilterType: (type: logFilterType) => set({filterType: type}),
  selectFilterOperation: (operation: filterOperations) =>
    set({filterOperation: operation}),
}))
