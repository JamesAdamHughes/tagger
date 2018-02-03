insert into tb_tag (tag_name)
		select
			'test4'
		where not exists(
			select
				pk_tag_id
			from tb_tag
			where tag_name = 'test4'
			limit 1
		);