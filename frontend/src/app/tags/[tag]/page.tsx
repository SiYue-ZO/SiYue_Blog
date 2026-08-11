import { redirect } from "next/navigation";

type Props = { params: Promise<{ tag: string }> };

export default async function LegacyTagPage({ params }: Props) {
  const { tag } = await params;
  redirect(`/blog?tag=${encodeURIComponent(tag)}`);
}
