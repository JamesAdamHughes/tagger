package tags

import "fmt"

func AddTag(songId string, tagId int64, tagName string, userId *int64) error {
	query :=  fmt.Sprintf(`
	insert into tb_user_song_tags (fk_user_id, fk_song_id, fk_tag_id)
	select 
		%d,
		$s,
		$d
	`, userId, songId, tagId)
}