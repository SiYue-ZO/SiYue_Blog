export default function Loading() {
  return <div className="page-shell py-16" aria-busy="true" aria-label="正在加载"><div className="h-5 w-24 animate-pulse rounded bg-muted" /><div className="mt-5 h-12 max-w-xl animate-pulse rounded-xl bg-muted" /><div className="mt-10 grid gap-5 md:grid-cols-3">{[0, 1, 2].map((item) => <div key={item} className="h-56 animate-pulse rounded-xl bg-muted" />)}</div></div>;
}
