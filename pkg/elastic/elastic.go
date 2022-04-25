package elastic

import (
	"context"
	"errors"
	"fing/log"
	"fing/pkg/db"
	"fing/pkg/entity/usr"
	"github.com/olivere/elastic/v7"
	"strconv"
	"time"
)

const (
	UserInfoIndex = "user_info"
)

type Elastic struct {
	Index string
}

func Index(index string) *Elastic {

	el := Elastic{}

	el.Index = index

	return &el
}

func (e *Elastic) DelUser(id ...string) {
	if db.EsClient == nil {
		return
	}

	ctx, cal := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cal()

	for _, v := range id {
		if _, err := db.EsClient.Delete().Index(e.Index).
			Id(v).Do(ctx); err != nil {
			log.Error("es删除用户ID:%v 信息失败:%v", v, err)
			return
		}
	}

}

func (e *Elastic) AddUser(user *usr.UserInfo) {
	if db.EsClient == nil {
		return
	}
	ctx, cal := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cal()

	_, err := db.EsClient.Index().Index(e.Index).
		Id(strconv.FormatInt(int64(user.ID), 10)).
		Refresh("wait_for").BodyJson(user).Do(ctx)
	if err != nil {
		log.Error("es添加用户ID:%v 信息失败:%v", user.ID, err)
		return
	}

}

func (e *Elastic) UpdateUser(id string, doc map[string]string) {
	if db.EsClient == nil {
		return
	}
	ctx, cal := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cal()

	_, err := db.EsClient.Update().Index(e.Index).
		Id(id).Doc(doc).Do(ctx)
	if err != nil {
		log.Error("es更新用户ID:%v 信息失败:%v", id, err)
		return
	}

}

func (e *Elastic) FindUserID(condition, keyword string, size int) ([]int64, error) {

	if db.EsClient == nil {
		return nil, errors.New("es is nil")
	}
	query := elastic.NewMatchPhraseQuery(condition, keyword)
	searchService := db.EsClient.Search().Index(e.Index).
		Query(query).From(0).Size(size)
	ctx, cal := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cal()

	resp, err := searchService.Do(ctx)
	if err != nil {
		return nil, err
	}

	if resp.TotalHits() == 0 {
		return nil, nil
	}

	users := make([]int64, 0, resp.TotalHits())
	for _, h := range resp.Hits.Hits {
		userID, _ := strconv.ParseInt(h.Id, 10, 64)
		if userID == 0 {
			continue
		}
		users = append(users, userID)
	}

	return users, nil
}
