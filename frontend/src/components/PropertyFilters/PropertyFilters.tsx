import React from 'react'
import { PropertyType, DealType } from '@/types/property'
import {
  Box,
  Button,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  SelectChangeEvent,
  TextField,
  Typography,
} from '@mui/material'
import { observer } from 'mobx-react-lite'
import { useStore } from '@/store'
import styles from './PropertyFilters.module.css'

export const PropertyFilters = observer(() => {
  const { browse } = useStore()

  const handlePropertyTypeChange = (event: SelectChangeEvent<PropertyType | ''>) => {
    const value = event.target.value as PropertyType | ''
    if (value === '') {
      browse.setFilter('propertyType', undefined)
    } else {
      browse.setFilter('propertyType', value as PropertyType)
    }
  }

  const handleDealTypeChange = (event: SelectChangeEvent<DealType | ''>) => {
    const value = event.target.value as DealType | ''
    if (value === '') {
      browse.setFilter('dealType', undefined)
    } else {
      browse.setFilter('dealType', value as DealType)
    }
  }

  const handleMinPriceChange = (event: any) => {
    const value = event.target.value ? Number(event.target.value) : undefined
    browse.setPriceRange(value, browse.filters.maxPrice)
  }

  const handleMaxPriceChange = (event: any) => {
    const value = event.target.value ? Number(event.target.value) : undefined
    browse.setPriceRange(browse.filters.minPrice, value)
  }

  const handleRoomsChange = (event: SelectChangeEvent<number | ''>) => {
    const value = event.target.value as number | ''
    browse.setRoomsCount(value === '' ? undefined : value)
  }

  const handleClearFilters = () => {
    browse.clearFilters()
  }

  return (
    <Box className={styles.filtersContainer}>
      <Typography variant='h6' className={styles.filtersTitle}>
        Фильтры
      </Typography>

      <Box className={styles.filtersRow}>
        <FormControl className={styles.filterControl} size='small'>
          <InputLabel id='property-type-label'>Тип недвижимости</InputLabel>
          <Select
            labelId='property-type-label'
            value={browse.filters.propertyType || ''}
            onChange={handlePropertyTypeChange}
            label='Тип недвижимости'
          >
            <MenuItem value=''>Все типы</MenuItem>
            <MenuItem value={PropertyType.APARTMENT}>Квартира</MenuItem>
            <MenuItem value={PropertyType.HOUSE}>Дом</MenuItem>
          </Select>
        </FormControl>

        <FormControl className={styles.filterControl} size='small'>
          <InputLabel id='deal-type-label'>Тип сделки</InputLabel>
          <Select
            labelId='deal-type-label'
            value={browse.filters.dealType || ''}
            onChange={handleDealTypeChange}
            label='Тип сделки'
          >
            <MenuItem value=''>Все типы</MenuItem>
            <MenuItem value={DealType.SALE}>Продажа</MenuItem>
            <MenuItem value={DealType.RENT}>Аренда</MenuItem>
          </Select>
        </FormControl>
      </Box>

      <Box className={styles.filtersRow}>
        <TextField
          className={styles.filterControl}
          label='Мин. цена'
          type='number'
          size='small'
          value={browse.filters.minPrice || ''}
          onChange={handleMinPriceChange}
        />
        <TextField
          className={styles.filterControl}
          label='Макс. цена'
          type='number'
          size='small'
          value={browse.filters.maxPrice || ''}
          onChange={handleMaxPriceChange}
        />
      </Box>

      <Box className={styles.filtersRow}>
        <FormControl className={styles.filterControl} size='small'>
          <InputLabel id='rooms-label'>Комнаты</InputLabel>
          <Select
            labelId='rooms-label'
            value={browse.filters.rooms || ''}
            onChange={handleRoomsChange}
            label='Комнаты'
          >
            <MenuItem value=''>Любое количество</MenuItem>
            <MenuItem value={1}>1</MenuItem>
            <MenuItem value={2}>2</MenuItem>
            <MenuItem value={3}>3</MenuItem>
            <MenuItem value={4}>4+</MenuItem>
          </Select>
        </FormControl>
      </Box>

      <Box className={styles.filtersActions}>
        <Button variant='outlined' onClick={handleClearFilters} className={styles.clearButton}>
          Сбросить фильтры
        </Button>
      </Box>
    </Box>
  )
})
