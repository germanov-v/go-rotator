package integrations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/germanov-v/go-rotator/internal/config"
	"github.com/germanov-v/go-rotator/internal/repository/postgres"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// const slot1 = "слот1";
const (
	slot1   = "slot1"
	slot2   = "slot2"
	slot3   = "slot3"
	bannerA = "bannerA"
	bannerB = "bannerB"
	bannerC = "bannerC"
	bannerD = "bannerD"
	bannerE = "bannerE"
	bannerF = "bannerF"
	bannerG = "bannerG"
	group1  = "group1"
	group2  = "group2"
	group3  = "group3"
)

func RunAll(ctx context.Context, cfg *config.Config) error {
	host := fmt.Sprintf("%s:%d", cfg.ServerConfig.ApiGatewayHost, cfg.ServerConfig.Port)

	repo, err := postgres.NewPostgresRepo(cfg.DbBaseConfig.ConnectionString)
	if err != nil {
		return fmt.Errorf("initi repo: %w", err)
	}

	if err := repo.AddGroup(ctx, group1); err != nil {
		return fmt.Errorf("add group: %w", err)
	}
	if err := repo.AddGroup(ctx, group2); err != nil {
		return fmt.Errorf("add group: %w", err)
	}
	if err := repo.AddGroup(ctx, group3); err != nil {
		return fmt.Errorf("add group: %w", err)
	}

	// http_handler.AddBannerHandler(repo)).Methods("POST")
	if err := testAddBanner(host, slot1, []string{bannerA, bannerB}); err != nil {
		return fmt.Errorf("testAddBanner failed: %w", err)
	}

	// http_handler.RotateBannerHandler(service)).Methods("GET")
	// Error expect 400 http_handler.RotateBannerHandler(service)).Methods("GET")
	err = testRotateMissingGroup(host, slot1)
	if err != nil {
		return fmt.Errorf("testRotateMissingGroup failed: %w", err)
	}
	//
	//// http_handler.RotateBannerHandler(service)).Methods("GET")
	//// Rotate  group + add click
	banner, err := testRotate(host, slot1, group1)
	if err != nil {
		return fmt.Errorf("testRotate failed: %w", err)
	}
	err = testRecordClick(host, slot1, banner, group1)
	if err != nil {
		return fmt.Errorf("testRecordClick failed: %w", err)
	}
	//

	//
	//// least once: Перебор всех: после большого количества показов, каждый баннер должен быть показан хотя один раз.

	err = testRotateAllBanners(host, slot2, []string{bannerD, bannerC, bannerE}, group2)
	if err != nil {
		return fmt.Errorf("testRotateAllBanners failed: %w", err)
	}
	//
	//// popular
	if err := testPopularBanner(host, slot3, []string{bannerF, bannerG}, group3); err != nil {
		return fmt.Errorf("testPopularBanner failed: %w", err)
	}

	//// http_handler.RemoveBannerHandler(repo)).Methods("DELETE")
	if err := testRemoveBanner(host, slot2, bannerD); err != nil {
		return fmt.Errorf("testRemoveBanner failed: %w", err)
	}

	return nil
}

func testAddBanner(host, slot string, banners []string) error {
	for _, b := range banners {
		// todo: в отдельный метод или ендпоинты разные и не стоит?
		url := fmt.Sprintf("%s/slots/%s/banners", host, slot)
		body := map[string]string{"banner_id": b}
		data, _ := json.Marshal(body)
		resp, err := http.Post(url, "application/json", bytes.NewReader(data))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			b, _ := ioutil.ReadAll(resp.Body)
			return fmt.Errorf("expected 201, got %d: %s", resp.StatusCode, string(b))
		}
	}
	return nil
}

func testRotateMissingGroup(host, slot string) error {
	url := fmt.Sprintf("%s/slots/%s/rotate", host, slot)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//defer func(Body io.ReadCloser) {
	//	err := Body.Close()
	//	if err != nil {
	//
	//	}
	//}(resp.Body)
	if resp.StatusCode != http.StatusBadRequest {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("expected 400, got %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func testRotate(host, slot, group string) (string, error) {
	url := fmt.Sprintf("%s/slots/%s/rotate?group=%s", host, slot, group)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		//panic("ALARM ERROR "+string(rune(resp.StatusCode)))
		return "", fmt.Errorf("===> expected 200, got %d: %s", resp.StatusCode, string(b))
	}
	var data struct {
		BannerId string `json:"banner_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	return data.BannerId, nil
}

func testRecordClick(host, slot, banner, group string) error {
	url := fmt.Sprintf("%s/slots/%s/stats/click", host, slot)
	body := map[string]string{"banner_id": banner, "group_id": group}
	data, _ := json.Marshal(body)

	resp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("====>>>> expected 200, got %d: %s", resp.StatusCode, string(b))
	}

	return nil
}

func testRemoveBanner(host, slot, banner string) error {
	url := fmt.Sprintf("%s/slots/%s/banners/%s", host, slot, banner)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("expected 200, got %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func testRotateAllBanners(host, slot string, banners []string, group string) error {
	err := testAddBanner(host, slot, banners)
	if err != nil {
		return err
	}
	counts := make(map[string]int)
	for i := 0; i < 100; i++ {
		banner, err := testRotate(host, slot, group)
		if err != nil {
			return err
		}
		//
		//if banner == bannerD {
		//	panic("DEMO TEST bannerD failed!")
		//}

		counts[banner]++
	}
	for _, b := range banners {
		if counts[b] == 0 {
			return fmt.Errorf("banner display %s lost", b)
		}
	}
	return nil
}

func testPopularBanner(host, slot string, banners []string, group string) error {
	if err := testAddBanner(host, slot, banners); err != nil {
		return err
	}
	// click 0 banner
	for i := 0; i < len(banners); i++ {
		banner, err := testRotate(host, slot, group)
		if err != nil {
			return err
		}
		if banner == banners[0] {
			if err := testRecordClick(host, slot, banner, group); err != nil {
				return err
			}
		}
	}

	// подкручиваем для нулевого стату
	for i := 0; i < 50; i++ {
		banner, err := testRotate(host, slot, group)
		if err != nil {
			return err
		}
		if banner == banners[0] {
			if err := testRecordClick(host, slot, banner, group); err != nil {
				return err
			}
		}
	}
	// TODO: ждем больше 50 получается? после подкрутки
	counts := make(map[string]int)
	for i := 0; i < 50; i++ {
		banner, err := testRotate(host, slot, group)
		if err != nil {
			return err
		}
		counts[banner]++
		// TODO: пока в один поток,
		time.Sleep(10 * time.Millisecond)
	}
	if counts[banners[0]] <= counts[banners[1]] {
		//return fmt.Errorf("counts[banners[0]] <= counts[banners[1]] IS TRUE  %d", banners[0])
		return fmt.Errorf("expected %s more popular: %d vs %d", banners[0], counts[banners[0]], counts[banners[1]])
	}
	return nil
}
