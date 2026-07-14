import z from "zod";

export const registerSchema = z.object({
  name: z
    .string("Name is required.")
    .trim()
    .min(3, "Name must be at least 3 characters long.")
    .max(60, "Name cannot exceed 60 characters."),

  email: z.email("Please enter a valid email address."),

  password: z
    .string("Password is required.")
    .min(8, "Password must be at least 8 characters long."),

  favorite_genres: z
    .array(
      z.object({
        genre_id: z.number("Genre ID is required."),
        name: z.string("Genre name is required."),
      }),
    )
    .min(1, "Please select at least one favorite genre."), // Optional: ensures they pick at least one
});

export type RegisterSchema = z.infer<typeof registerSchema>;
