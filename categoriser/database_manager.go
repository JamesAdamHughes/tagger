package categoriser

import (
	"tagger/server/database"
)

type AddSongTagRequest struct {
	SongId string
	UserId string
	TagId int64
	TagName string
}

type GetSongTagRequest struct {
	SongId string
	UserId string
}


func AddSongTag(request *AddSongTagRequest) (err error) {

	// add the tag to tb_tag if needed + get tag id
	tag, err := AddTag(request.TagName)
	if err != nil {
		return err
	}

	// Add id to request
	request.TagId = tag.ID

	// Insert tag
	query := `
	insert into tb_user_song_tags (fk_user_id, fk_song_id, fk_tag_id)
	select 
		:user_id,
		:song_id,
		:tag_id
	where not exists (
		select 
			pk_user_song_tag_id
		from tb_user_song_tags
		where 1=1 
			and fk_song_id = :song_id
			and fk_tag_id = :tag_id
		
	)
	`

	_, err = database.Exec(query, map[string]interface{}{
		"song_id": request.SongId,
		"tag_id": request.TagId,
		"user_id": request.UserId,
	})

	return
}

func AddTag(tagName string) (tag *Tag, err error) {
	query := `
		insert into tb_tag (tag_name)
		select 
			:tag_name 	
		where not exists(
			select  
				pk_tag_id 
			from tb_tag 	
			where 1=1
				and tag_name = :tag_name 	
			limit 1
		)
	`

	selectQuery := `
		select 
			pk_tag_id,
			tag_name
		from tb_tag 
		where 1=1
			and tag_name = :tag_name
		limit 1 
	`

	// Check insert the tag if it doesn't exist
	_, err = database.Exec(query, map[string]interface{}{
		"tag_name": tagName,
	})

	if err != nil {
		return nil, err
	}

	// Get the tag id from name
	rows, err := database.Select(selectQuery, map[string]interface{}{
		"tag_name": tagName,
	})
	if err != nil {
		return nil, err
	}

	var t Tag
	for rows.Next() {
		err = rows.Scan(&t.ID, &t.Name)
		if err != nil {
			return nil, err
		}
	}

	tag = &t
	return
}

func GetSongTags(request GetSongTagRequest) (tags []Tag, err error) {
	query := `
		select
			fk_tag_id,
			tag_name
		from tb_user_song_tags tbust 
		join tb_tag tbt
			on tbt.pk_tag_id = tbust.fk_tag_id
		where 1=1
			and tbust.fk_song_id = :song_id
			and tbust.fk_user_id = :user_id
			
	`

	rows, err := database.Select(query,map[string]interface{}{
		"song_id": request.SongId,
		"user_id": request.UserId,
	})
	if err != nil {
		return nil, err
	}


	for rows.Next() {
		var t Tag
		err := rows.Scan(&t.ID, &t.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	return
}