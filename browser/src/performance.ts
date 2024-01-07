import * as LH from "lighthouse/types/lh";
import {Result} from 'lighthouse'

export interface Performance {
  title: string
  score: number | null
  refs: Array<{name: string, weight: number}>
}

/**
 * JSON形式でMetricsを取得
 * @param result
 */
export function convertJson(result: LH.RunnerResult): Performance[] {
  const performances: Performance[] = [];
  Object.keys(result.lhr.categories).forEach(key => {
    const r: Result.Category = result.lhr.categories[key];

    const p: Performance = {
      title: r.title,
      score: r.score,
      refs: r.auditRefs.map((ref) => {
        return { name: ref.id, weight: ref.weight }
      }),
    }
    performances.push(p);
  });

  return performances;
}
