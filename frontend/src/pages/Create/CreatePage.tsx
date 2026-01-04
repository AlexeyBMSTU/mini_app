import React, { useState } from 'react'
import { motion } from 'motion/react'
import {
  Button,
  FormControl,
  FormControlLabel,
  FormLabel,
  MenuItem,
  Select,
  Switch,
  TextField,
  Typography,
  Box,
  CircularProgress,
} from '@mui/material'
import { PropertyType, DealType } from '@/types/property'
import apiService from '@/services/apiService'
import styles from '@/pages/Create/CreatePage.module.css'
import { PurePage } from '@/components/PurePage/PurePage'

export const CreatePage = () => {
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    type: PropertyType.APARTMENT,
    dealType: DealType.SALE,
    price: '',
    area: '',
    rooms: '',
    floor: '',
    totalFloors: '',
    yearBuilt: '',
    address: '',
    city: '',
    district: '',
    features: {
      hasBalcony: false,
      hasParking: false,
      hasElevator: false,
      hasFurniture: false,
      hasKitchen: false,
      hasInternet: false,
      hasAirConditioning: false,
      hasWashingMachine: false,
      hasDishwasher: false,
      hasTV: false,
      hasRefrigerator: false,
    },
  })

  const [errors, setErrors] = useState<Record<string, string>>({})
  const [isSubmitting, setIsSubmitting] = useState(false)
  const requiredError = 'Обязательное поле'
  const incorrectField = 'Неккоректное поле'

  const validateForm = () => {
    const newErrors: Record<string, string> = {}

    if (!formData.title.trim()) {
      newErrors.title = requiredError
    }

    if (!formData.description.trim()) {
      newErrors.description = requiredError
    }

    if (!formData.price) {
      newErrors.price = requiredError
    } else if (isNaN(Number(formData.price)) || Number(formData.price) <= 0) {
      newErrors.price = incorrectField
    }

    if (!formData.area) {
      newErrors.area = requiredError
    } else if (isNaN(Number(formData.area)) || Number(formData.area) <= 0) {
      newErrors.area = incorrectField
    }

    if (!formData.rooms) {
      newErrors.rooms = requiredError
    } else if (isNaN(Number(formData.rooms)) || Number(formData.rooms) <= 0) {
      newErrors.rooms = incorrectField
    }

    if (!formData.city.trim()) {
      newErrors.city = requiredError
    }

    if (!formData.address.trim()) {
      newErrors.address = requiredError
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    if (name) {
      setFormData(prev => ({
        ...prev,
        [name]: value,
      }))

      if (errors[name]) {
        setErrors(prev => ({
          ...prev,
          [name]: '',
        }))
      }
    }
  }

  const handleSelectChange = (name: string) => (e: any) => {
    setFormData(prev => ({
      ...prev,
      [name]: e.target.value,
    }))

    if (errors[name]) {
      setErrors(prev => ({
        ...prev,
        [name]: '',
      }))
    }
  }

  const handleFeatureChange =
    (featureName: string) => (event: React.ChangeEvent<HTMLInputElement>) => {
      setFormData(prev => ({
        ...prev,
        features: {
          ...prev.features,
          [featureName]: event.target.checked,
        },
      }))
    }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (validateForm()) {
      setIsSubmitting(true)
      try {
        const propertyData = {
          title: formData.title,
          description: formData.description,
          type: formData.type,
          dealType: formData.dealType,
          price: Number(formData.price),
          area: Number(formData.area),
          rooms: Number(formData.rooms),
          floor: formData.floor ? Number(formData.floor) : undefined,
          totalFloors: formData.totalFloors ? Number(formData.totalFloors) : undefined,
          yearBuilt: formData.yearBuilt ? Number(formData.yearBuilt) : undefined,
          address: formData.address,
          city: formData.city,
          district: formData.district || undefined,
          features: formData.features,
        }

        await apiService.createProperty(propertyData)
      } catch (error) {
        console.error('Error creating property:', error)
      } finally {
        setIsSubmitting(false)
      }
    }
  }

  return (
    <PurePage>
      <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} transition={{ duration: 0.5 }}>
        <form onSubmit={handleSubmit}>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
            <Box>
              <TextField
                fullWidth
                label='Заголовок объявления*'
                name='title'
                value={formData.title}
                onChange={handleChange}
                error={!!errors.title}
                helperText={errors.title}
              />
            </Box>

            <Box>
              <TextField
                fullWidth
                label='Описание*'
                name='description'
                value={formData.description}
                onChange={handleChange}
                multiline
                rows={4}
                error={!!errors.description}
                helperText={errors.description}
              />
            </Box>

            <Box sx={{ display: 'flex', gap: 2, flexDirection: { xs: 'column', sm: 'row' } }}>
              <FormControl fullWidth>
                <FormLabel>Тип недвижимости</FormLabel>
                <Select name='type' value={formData.type} onChange={handleSelectChange('type')}>
                  <MenuItem value={PropertyType.APARTMENT}>Квартира</MenuItem>
                  <MenuItem value={PropertyType.HOUSE}>Дом</MenuItem>
                </Select>
              </FormControl>

              <FormControl fullWidth>
                <FormLabel>Тип сделки</FormLabel>
                <Select
                  name='dealType'
                  value={formData.dealType}
                  onChange={handleSelectChange('dealType')}
                >
                  <MenuItem value={DealType.SALE}>Продажа</MenuItem>
                  <MenuItem value={DealType.RENT}>Аренда</MenuItem>
                </Select>
              </FormControl>
            </Box>

            <Box sx={{ display: 'flex', gap: 2, flexDirection: { xs: 'column', sm: 'row' } }}>
              <TextField
                fullWidth
                label='Цена*'
                name='price'
                type='number'
                value={formData.price}
                onChange={handleChange}
                error={!!errors.price}
                helperText={errors.price}
              />

              <TextField
                fullWidth
                label='Площадь (м²)*'
                name='area'
                type='number'
                value={formData.area}
                onChange={handleChange}
                error={!!errors.area}
                helperText={errors.area}
              />
            </Box>

            <Box sx={{ display: 'flex', gap: 2, flexDirection: { xs: 'column', sm: 'row' } }}>
              <TextField
                fullWidth
                label='Количество комнат*'
                name='rooms'
                type='number'
                value={formData.rooms}
                onChange={handleChange}
                error={!!errors.rooms}
                helperText={errors.rooms}
              />

              <TextField
                fullWidth
                label='Этаж'
                name='floor'
                type='number'
                value={formData.floor}
                onChange={handleChange}
              />

              <TextField
                fullWidth
                label='Всего этажей'
                name='totalFloors'
                type='number'
                value={formData.totalFloors}
                onChange={handleChange}
              />
            </Box>

            <Box sx={{ display: 'flex', gap: 2, flexDirection: { xs: 'column', sm: 'row' } }}>
              <TextField
                fullWidth
                label='Год постройки'
                name='yearBuilt'
                type='number'
                value={formData.yearBuilt}
                onChange={handleChange}
              />

              <TextField
                fullWidth
                label='Город*'
                name='city'
                value={formData.city}
                onChange={handleChange}
                error={!!errors.city}
                helperText={errors.city}
              />
            </Box>

            <Box>
              <TextField
                fullWidth
                label='Адрес*'
                name='address'
                value={formData.address}
                onChange={handleChange}
                error={!!errors.address}
                helperText={errors.address}
              />
            </Box>

            <Box>
              <TextField
                fullWidth
                label='Район'
                name='district'
                value={formData.district}
                onChange={handleChange}
              />
            </Box>

            <Box>
              <Typography variant='h6'>Дополнительные особенности</Typography>
              <Box className={styles.featuresGrid}>
                <Box className={styles.featureItem}>
                  <FormControlLabel
                    control={
                      <Switch
                        checked={formData.features.hasBalcony}
                        onChange={handleFeatureChange('hasBalcony')}
                        name='hasBalcony'
                      />
                    }
                    label='Балкон'
                  />
                </Box>
                <Box className={styles.featureItem}>
                  <FormControlLabel
                    control={
                      <Switch
                        checked={formData.features.hasParking}
                        onChange={handleFeatureChange('hasParking')}
                        name='hasParking'
                      />
                    }
                    label='Парковка'
                  />
                </Box>
                <Box className={styles.featureItem}>
                  <FormControlLabel
                    control={
                      <Switch
                        checked={formData.features.hasElevator}
                        onChange={handleFeatureChange('hasElevator')}
                        name='hasElevator'
                      />
                    }
                    label='Лифт'
                  />
                </Box>
                <Box className={styles.featureItem}>
                  <FormControlLabel
                    control={
                      <Switch
                        checked={formData.features.hasFurniture}
                        onChange={handleFeatureChange('hasFurniture')}
                        name='hasFurniture'
                      />
                    }
                    label='Мебель'
                  />
                </Box>
                <Box className={styles.featureItem}>
                  <FormControlLabel
                    control={
                      <Switch
                        checked={formData.features.hasKitchen}
                        onChange={handleFeatureChange('hasKitchen')}
                        name='hasKitchen'
                      />
                    }
                    label='Кухня'
                  />
                </Box>
                <Box className={styles.featureItem}>
                  <FormControlLabel
                    control={
                      <Switch
                        checked={formData.features.hasInternet}
                        onChange={handleFeatureChange('hasInternet')}
                        name='hasInternet'
                      />
                    }
                    label='Интернет'
                  />
                </Box>
                <Box className={styles.featureItem}>
                  <FormControlLabel
                    control={
                      <Switch
                        checked={formData.features.hasAirConditioning}
                        onChange={handleFeatureChange('hasAirConditioning')}
                        name='hasAirConditioning'
                      />
                    }
                    label='Кондиционер'
                  />
                </Box>
                <Box className={styles.featureItem}>
                  <FormControlLabel
                    control={
                      <Switch
                        checked={formData.features.hasWashingMachine}
                        onChange={handleFeatureChange('hasWashingMachine')}
                        name='hasWashingMachine'
                      />
                    }
                    label='Стиральная машина'
                  />
                </Box>
                <Box className={styles.featureItem}>
                  <FormControlLabel
                    control={
                      <Switch
                        checked={formData.features.hasDishwasher}
                        onChange={handleFeatureChange('hasDishwasher')}
                        name='hasDishwasher'
                      />
                    }
                    label='Посудомоечная машина'
                  />
                </Box>
                <Box className={styles.featureItem}>
                  <FormControlLabel
                    control={
                      <Switch
                        checked={formData.features.hasTV}
                        onChange={handleFeatureChange('hasTV')}
                        name='hasTV'
                      />
                    }
                    label='Телевизор'
                  />
                </Box>
                <Box className={styles.featureItem}>
                  <FormControlLabel
                    control={
                      <Switch
                        checked={formData.features.hasRefrigerator}
                        onChange={handleFeatureChange('hasRefrigerator')}
                        name='hasRefrigerator'
                      />
                    }
                    label='Холодильник'
                  />
                </Box>
              </Box>
            </Box>

            <Box>
              <motion.div whileHover={{ scale: 1.02 }} whileTap={{ scale: 0.98 }}>
                <Button
                  type='submit'
                  variant='contained'
                  color='primary'
                  fullWidth
                  size='large'
                  disabled={isSubmitting}
                  startIcon={
                    isSubmitting ? <CircularProgress size={20} color='inherit' /> : undefined
                  }
                >
                  {isSubmitting ? 'Создание...' : 'Создать объявление'}
                </Button>
              </motion.div>
            </Box>
          </Box>
        </form>
      </motion.div>
    </PurePage>
  )
}
