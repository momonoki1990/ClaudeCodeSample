export type Category = {
  id: number
  name: string
  is_system: boolean
}

export type Todo = {
  id: number
  text: string
  done: boolean
  category_id: number | null
  category: Category | null
  position: number
}
