import { z } from 'zod'

export const AttachmentSchema = z.object({
  id: z.string(),
  name: z.string(),
  url: z.string().optional(),
  type: z.string().optional(),
  description: z.string().optional(),
  display: z.array(z.string()).optional(),
})

export const ProductLanguageSchema = z.object({
  id: z.string(),
  languageCode: z.string(),
  name: z.string(),
  description: z.string().optional(),
  attachment: z.array(AttachmentSchema).optional(),
})

const Price = z.object({
  value: z.number().optional(),
  unit: z.string().optional(),
})

export const ProductPriceLanguageSchema = z.object({
  id: z.string(),
  languageCode: z.string(),
  price: z.number().optional(),
  name: z.string().optional(),
  description: z.string().optional(),
  unit: z.string().optional(),
  vat: Price.optional(),
  tax: Price.optional(),
})

export const ProductPriceSchema = z.object({
  id: z.string(),
  value: z.number().optional(),
  unit: z.string().optional(),
  vat: Price.optional(),
  tax: Price.optional(),
  language: z.array(ProductPriceLanguageSchema).optional(),
})

export const ProductSchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string().optional(),
  price: ProductPriceSchema.optional(),
  language: z.array(ProductLanguageSchema).optional(),
})

export type IProductSchema = z.infer<typeof ProductSchema>

const Product = z.object({
  name: z.string(),
  description: z.string().optional(),
  price: z.number().optional(),
  unit: z.string().optional(),
  vat: z.number().optional(),
  tax: z.number().optional(),
})

export const ProductBodySchema = z.object({
  name: z.string(),
  description: z.string().optional(),
  price: z.number().optional(),
  unit: z.string().optional(),
  vat: z.number().optional(),
  tax: z.number().optional(),
  th: Product.optional(),
  en: Product.optional(),
})

export type IProductBody = z.infer<typeof ProductBodySchema> & {
  id?: string
}

export type IProduct = z.infer<typeof ProductSchema> & {
  href: string
}
