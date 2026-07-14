import z from "zod";

export const registerSchema = z.object({
  name: z.string().min(3).max(60),
  email: z.email(),
  password: z.string().min(8),
  favorite_genres: z.array(
    z.object({ genre_id: z.number(), name: z.string() }),
  ),
});

export type RegisterSchema = z.infer<typeof registerSchema>;
