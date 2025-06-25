-- マッチングアプリ用初期データ投入スクリプト
-- 実行前に 001_create_matching_app_tables.sql を実行してください

-- タグマスタデータ
INSERT INTO tags (name, created_at, updated_at) VALUES
('プログラミング', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('旅行', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('料理', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('映画', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('音楽', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('スポーツ', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('読書', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('ゲーム', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('アート', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('写真', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (name) DO NOTHING;

-- ユーザーテストデータ
INSERT INTO users (gmail, name, created_at, icon_url) VALUES
('taro.yamada@gmail.com', '山田太郎', CURRENT_TIMESTAMP, 'https://example.com/icons/taro.jpg'),
('hanako.sato@gmail.com', '佐藤花子', CURRENT_TIMESTAMP, 'https://example.com/icons/hanako.jpg'),
('kenji.tanaka@gmail.com', '田中健二', CURRENT_TIMESTAMP, 'https://example.com/icons/kenji.jpg'),
('yuki.watanabe@gmail.com', '渡辺雪', CURRENT_TIMESTAMP, 'https://example.com/icons/yuki.jpg'),
('hiroshi.ito@gmail.com', '伊藤寛', CURRENT_TIMESTAMP, 'https://example.com/icons/hiroshi.jpg'),
('ai.kobayashi@gmail.com', '小林愛', CURRENT_TIMESTAMP, 'https://example.com/icons/ai.jpg'),
('daisuke.kato@gmail.com', '加藤大輔', CURRENT_TIMESTAMP, 'https://example.com/icons/daisuke.jpg'),
('miki.yoshida@gmail.com', '吉田美紀', CURRENT_TIMESTAMP, 'https://example.com/icons/miki.jpg'),
('takeshi.nakamura@gmail.com', '中村武', CURRENT_TIMESTAMP, NULL),
('rina.hayashi@gmail.com', '林里奈', CURRENT_TIMESTAMP, NULL)
ON CONFLICT (gmail) DO NOTHING;

-- プロフィールテストデータ
-- ユーザーIDを動的に取得してプロフィールを作成
WITH user_tag_mapping AS (
  SELECT u.id as user_id, t.id as tag_id, u.name as user_name
  FROM users u
  CROSS JOIN tags t
  WHERE 
    (u.name = '山田太郎' AND t.name = 'プログラミング') OR
    (u.name = '佐藤花子' AND t.name = '旅行') OR
    (u.name = '田中健二' AND t.name = 'スポーツ') OR
    (u.name = '渡辺雪' AND t.name = '音楽') OR
    (u.name = '伊藤寛' AND t.name = '映画') OR
    (u.name = '小林愛' AND t.name = '料理') OR
    (u.name = '加藤大輔' AND t.name = 'ゲーム') OR
    (u.name = '吉田美紀' AND t.name = 'アート') OR
    (u.name = '中村武' AND t.name = '読書') OR
    (u.name = '林里奈' AND t.name = '写真')
)
INSERT INTO profiles (user_id, bio, tag_id, created_at, updated_at)
SELECT 
  utm.user_id,
  CASE 
    WHEN utm.user_name = '山田太郎' THEN 'フルスタックエンジニアです。React、Go、PostgreSQLが得意です。一緒にプロジェクトを作りませんか？'
    WHEN utm.user_name = '佐藤花子' THEN '旅行が大好きです！これまで30カ国以上を訪れました。新しい文化に触れることが楽しみです。'
    WHEN utm.user_name = '田中健二' THEN 'サッカーを20年続けています。週末はフットサルチームで汗を流しています。'
    WHEN utm.user_name = '渡辺雪' THEN 'ピアノとギターを演奏します。音楽を通じて人とのつながりを大切にしています。'
    WHEN utm.user_name = '伊藤寛' THEN '映画鑑賞が趣味です。特にSF映画が好きで、年間200本以上観ています。'
    WHEN utm.user_name = '小林愛' THEN '料理研究家として活動しています。健康的で美味しい料理を作ることが生きがいです。'
    WHEN utm.user_name = '加藤大輔' THEN 'ゲーム開発者です。インディーゲームから大型タイトルまで幅広く手がけています。'
    WHEN utm.user_name = '吉田美紀' THEN 'グラフィックデザイナーです。アートとデザインの境界を探求しています。'
    WHEN utm.user_name = '中村武' THEN '月20冊以上読書します。特に哲学書と小説が好きです。'
    WHEN utm.user_name = '林里奈' THEN 'フォトグラファーです。自然と人物写真を中心に撮影しています。'
  END,
  utm.tag_id,
  CURRENT_TIMESTAMP,
  CURRENT_TIMESTAMP
FROM user_tag_mapping utm
ON CONFLICT (user_id) DO NOTHING;

-- ブックマークテストデータ
WITH user_pairs AS (
  SELECT 
    u1.id as user_id,
    u2.id as bookmarked_user_id
  FROM users u1
  CROSS JOIN users u2
  WHERE u1.id != u2.id
    AND (
      (u1.name = '山田太郎' AND u2.name IN ('佐藤花子', '田中健二', '小林愛')) OR
      (u1.name = '佐藤花子' AND u2.name IN ('山田太郎', '渡辺雪', '林里奈')) OR
      (u1.name = '田中健二' AND u2.name IN ('山田太郎', '伊藤寛', '加藤大輔')) OR
      (u1.name = '渡辺雪' AND u2.name IN ('佐藤花子', '吉田美紀')) OR
      (u1.name = '伊藤寛' AND u2.name IN ('田中健二', '中村武'))
    )
)
INSERT INTO bookmarks (user_id, bookmarked_user_id, created_at)
SELECT user_id, bookmarked_user_id, CURRENT_TIMESTAMP
FROM user_pairs
ON CONFLICT (user_id, bookmarked_user_id) DO NOTHING;

-- マッチングテストデータ
WITH user_matchings AS (
  SELECT 
    u.id as user_id,
    u.name as user_name
  FROM users u
  WHERE u.name IN ('山田太郎', '佐藤花子', '田中健二', '渡辺雪', '小林愛')
)
INSERT INTO matchings (user_id, notify_id, content, created_at, updated_at)
SELECT 
  um.user_id,
  (RANDOM() * 1000)::INTEGER as notify_id,
  CASE 
    WHEN um.user_name = '山田太郎' THEN '佐藤花子さんとマッチしました！共通の趣味について話しましょう。'
    WHEN um.user_name = '佐藤花子' THEN '山田太郎さんとマッチしました！プログラミングについて教えてください。'
    WHEN um.user_name = '田中健二' THEN '新しいスポーツ仲間を探しています。一緒に運動しませんか？'
    WHEN um.user_name = '渡辺雪' THEN '音楽好きの方とマッチしました！セッションしませんか？'
    WHEN um.user_name = '小林愛' THEN '料理好きの方々とマッチしました。レシピ交換しましょう！'
  END,
  CURRENT_TIMESTAMP - INTERVAL '1 day' * (RANDOM() * 7),
  CURRENT_TIMESTAMP
FROM user_matchings um
ON CONFLICT DO NOTHING;

-- データ確認用クエリ（コメントアウト）
/*
-- 作成されたデータの確認
SELECT 'Users' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'Tags' as table_name, COUNT(*) as count FROM tags
UNION ALL
SELECT 'Profiles' as table_name, COUNT(*) as count FROM profiles
UNION ALL
SELECT 'Bookmarks' as table_name, COUNT(*) as count FROM bookmarks
UNION ALL
SELECT 'Matchings' as table_name, COUNT(*) as count FROM matchings;

-- ユーザーとプロフィールの関連確認
SELECT 
  u.name as user_name,
  u.gmail,
  p.bio,
  t.name as tag_name
FROM users u
LEFT JOIN profiles p ON u.id = p.user_id
LEFT JOIN tags t ON p.tag_id = t.id
ORDER BY u.id;

-- ブックマーク関係の確認
SELECT 
  u1.name as user_name,
  u2.name as bookmarked_user_name,
  b.created_at
FROM bookmarks b
JOIN users u1 ON b.user_id = u1.id
JOIN users u2 ON b.bookmarked_user_id = u2.id
ORDER BY b.created_at;
*/