package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// ============================================================================
// CONFIGURATION
// ============================================================================

const (
	GraphQLEndpoint = "https://bff-page.kakao.com/graphql"

	// Thay Cookie mới của bạn vào đây
	CookieString = "_kpdid=193cdcc72d6d4125917f8fa486d1337f; _kpiid=26fd12766e8056921aa387ec10f166af; _kahai=ba4be1e91c915262b60a6402f5e1e321446b1a1739f6a79ad21e84f1393c8927; _kawask=03074906-0b6a-4b25-a2fa-a5ff94e956d3; _kau=3f523a174a652d4235aca16fd27277b822346dcb30439688b66c44b42ed6b1f455db8ded9b0a76c6defda3d52f3c3be77c0962a4b82775b26fe28cac561b97b5c068dbbdbe642261e352c83f9d9dffb59c3cb00295a6f514c53cac54041955f59929adcbd676b62c190dee2f31dde821e8b4c7f630f7dce7cbd94ddf3639323730393434353636353534353037343432343330383837313036343432f7d9933049124eb1f43a63ef8a797e21; _kawlt=bMlacF1u-YmEOt7pPmKt6RShvZs9uBuI6Or4XpOIIcycSFvGVX_YqeJ7Js-p2JQxDJBGukqxqnp-3jK2HY2AMIgBFq75kQlfH-LTOenFXQjPl-oeseU82EPNLPo_U-Zo; _kawltea=1764935369; _karmt=EMnRgDJz-ZDs4gKdhvdhHUeDYqZIg5SSd-h0Y-w5bwE8d7wXRP42dUJw4ys6y5Q_; _karmtea=1767462569; _kaslt=EaosIf61j+DpRVnPDuqZmgnRJ1dJTQsvuTpBlYqpfAxJ2+9HCyAORx8k+zMMywPgLwyaJgBSiSdMB9zmD5NbBA==; _kpawbat_e=8pDSy1%2Fb0KhSfA5AxsB56FKP%2B1uHZVh7r74yxeibkAUUuC5F86ADw5OANejxyRO9ZS2eaA4GkUEkQpKRyERBR86SAjSdpAcbQKBG%2FJuTu0w%3D; _kpwtkn=eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNIn0.._adSp9e8oto0_lkg.aH_opXKck2d--Oq3g7hXwdcxtTQWITa0xoEL5UHaIYWpGj0fBk7qHtZre_1eaiisThjJJeA2DPcn6R5ni2_84vf-Ll9Q7RD4D8Y_GWwf2A3jOdYp91_kCRZlBv8zjomNX5HOsjGTeAZyK_uZ8up0wvGUIzv-xzrR_cgsXCOkyQslW8XQtYsfEjegDrCxnGOfn5cTAhkllUi5jlMlz-5V-SKtvNBQHbHJ_oxkoPZI5cLO2Oh6L-RWgRK9eu-aLHRXHGGSGkOm5EzWsUpN5-oHpF5FlkdnABX3_Y_FVkAx3rK_vDtIfqSG1vcUHfAHI5NHes8dHLHaGge2GpbIjOa5Z50z7VSyeuSllLdN7mXtD7M3YnN244szQlsD-yvMDcvCU8WGbfFDzl8dnxPXn8NkSKCCA4gXFkKthI7Vxgc-gNKg7B9muYk2Huslcdb6GAdRg-N7Tbip6t5s6vsZXMmi3GMhJL-3teC5S-wuNTcvVE-uasJgAKyJda04dDUcgAHXOKuyVoIVKPOh3efMhIR6LjqeUUK8iXj8_kpkjuyqEYwX4PjdEHRYwQCzMrsyKyH0fczsMkSIcLCbDWerigc1gB8FWuvGUL9hxlNFh-lsabqIqcxqpyc5iT9sm3QYdBrgiYmJn8TNOmc4_cc5SvaAeQOqzq1IgwgY2u962e8r-thzWd_SQPAwt_cTDA8VUarxLLPPnl0bq5WFSZWx3xsOmZdnRmCEZ4HB0twugaoL9k2CwZJEpAY_dVWA4vUuz0dwOtxukShPhOLtmVbpNOU1LpAxuGDIDEERsjd2vceNvcOcaMpsc_kwgSE4Ng3R376rvVx0yGlH3VUFBqOt0pmpumqnf17Io7_yNz4yr2dcW9c0NfxHYn5Qf-3oLRufh-nZK1dK6XI2iy0wckfeczM4xRoBKCf9bRvecnMGGaP67zvv7OsmN2E7MoapjwlTtwOTMkhn5IjSJhhC65-08AdtLdDvjAbhmyWRqfycHw0Eqg0koqlgFzmELZfZ_-oG7Ks0qi3BD-cUGqhN3XZSTqoHRg5_RQV5qaGo_1sz1ZyTc_cjS0i5uT8R2DakUTRsA9bCPJYpAKMWTRtX_h9oRwRPKoNkUpPCYbw7AS0Ui9VrsTyE7lgymqdbcu2bgriI6s4B6EM-j-ZwzI2h2o6AoV932B0OUuCYeGj1WaQVjEMKonh3B_U1bTbnR5u6qq82q2ZpsMAWwbInv7-zIj0IU-fmQKRB4Nw4arvF4pap2nAPxP5LvWL1DN85-GKS-pOea0ujcaxgQflkLT_1pNw3PABk42x7J_JRcALf-3jwLr_GUoZ7FtKarKYLsYvmeNdHi9Gri5uqlXfj3OKpbCjjpr37XLs7JjiEF0opHxpsQiCfRsurGp5AD-7gIyr9CA_jc_C_DDROYsVcp_H99l4G2thj-0djp-TKjSM4vO8MRHFK6PRaSZwbO_K-q-pCfZQW7x--xH3MC28cCJuXHoBxhPXKX3IsQQNF0Xvy_bJAjZj2mxDEJ3uxsYTpPwRvSdzzweZzPgaqwDpovL2wEHlF43FKXCN8XJgsgPoOlgphtZZbjKxeLI4GBX-GCknfeYZBeNBBcfRzzc5dAszIExqWHt8C_Jr2dqpyXzBvPbOarTVs5HT881BYmhu7_0nkSHXHtu1JcNOZvGTatHT7DFH8g02Qd-OkZ3yqZL4KUQ2GfMsQLdyx-f2lRmeb8VGxL1wJcG6rPqQuZf59MUGXkiaNSJNbv6shwv50NysEtxeJu-mD-0D0mHfFHf5dtcJbu_w-bOUCdX0LxQZRGRicP_OS-mcvcNd5g-crsjPJmFzG6_i4-n0l6KplXul8PezHnpfMO-2afE1mv7dTBYjBJNJHl3fuNbpzBjVSTnb9uUQd_lJwYHtqod5Vl4HlVCXfa_lMz0epn_yWsyE-alFxuZd4UyIuKnfYU1zQ69851wsSgBCzviT_AOgh3THVJKCZjPWOZjM_dfsK4FknMbVnLgi64-VZfE22I74GKTW8WXNTEtwL8zmKkYQGNsO3PME9NjJb93rs3iZnzYyN5l_JOwcXe80M0__0tVXXFYu8jImBinQKNHiMPvLEbFPiHj0HhW5MP2MydDGtY3AyX-DXyyqdbkq0qx0M9cHp3bLwqP6Ar1kGps6yifqq65BlVS_EXbhnid-ndxzJvUeu-fu48hbN82uB67EwMHf-bFofHP-Z5_n3kMCoTpThtRw2AxS2A2A4dNttt2TeBI-zILItzevvN0UKI-Rvj8XEPNJUoIU5393y9nPCzHESkdZoNmXpwwbfwvg.rDjFii5fiYRJfBkC1xgdiQ; _T_ANO=Zw/bxT4wI+Kd1FDDKJ3UJRQsaKoThbjt0ICL9cEc3wnUockARZDkb4WYxkYq0a+zDSuKWh18gnbwpkdvWfW+Hb8THP7Ra6aMpepyV+j9kstQZsKptAKBzfkMncexAiJuWEitZFYi92O3Eijhry6S2Lh1XCBWZN9M58g6VYw/xf7HmPTwyQYFRfTkWkOmVkBp/lqaDPUSDT+xk824X6gUcP2al1BvY92E02zTNlHtNII/ekkMHfQdzfZMNmV65p/Uy5juQjycuf2llbSprnYavIGZLqIK7ZdYLYPNLF+IoQXMpZcfHn0+XywsHKXXEh9F3lv3DcL0qrxtMZUNrebepg==" // rút gọn
)

// ============================================================================
// QUERY (dùng payload bạn cung cấp — giữ nguyên các type Long!)
// ============================================================================

const ViewerQuery = `
query viewerInfo($seriesId: Long!, $productId: Long!) {
  viewerInfo(seriesId: $seriesId, productId: $productId) {
    item {
      ...SingleFragment
    }
    seriesItem {
      ...SeriesFragment
    }
    prevItem {
      ...NearItemFragment
    }
    nextItem {
      ...NearItemFragment
    }
    viewerData {
      ...TextViewerData
      ...TalkViewerData
      ...ImageViewerData
      ...VodViewerData
    }
    displayAd {
      ...DisplayAd
    }
  }
}

fragment SingleFragment on Single {
  id
  productId
  seriesId
  title
  thumbnail
  badge
  isFree
  ageGrade
  state
  slideType
  lastReleasedDate
  size
  pageCount
  isHidden
  remainText
  isWaitfreeBlocked
  saleState
  series {
    ...SeriesFragment
  }
  serviceProperty {
    ...ServicePropertyFragment
  }
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
  discountRate
  discountRateText
  isShortsDrama
}

fragment SeriesFragment on Series {
  id
  seriesId
  title
  thumbnail
  landThumbnail
  categoryUid
  lang
  category
  categoryType
  subcategoryUid
  subcategory
  badge
  isAllFree
  isWaitfree
  ageGrade
  state
  onIssue
  authors
  description
  pubPeriod
  freeSlideCount
  lastSlideAddedDate
  waitfreeBlockCount
  waitfreePeriodByMinute
  bm
  saleState
  startSaleDt
  saleMethod
  discountRate
  discountRateText
  serviceProperty {
    ...ServicePropertyFragment
  }
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
  translateProperty {
    ...TranslatePropertyFragment
  }
}

fragment ServicePropertyFragment on ServiceProperty {
  viewCount
  readCount
  ratingCount
  ratingSum
  commentCount
  pageContinue {
    ...ContinueInfoFragment
  }
  todayGift {
    ...TodayGift
  }
  preview {
    ...PreviewFragment
  }
  waitfreeTicket {
    ...WaitfreeTicketFragment
  }
  isAlarmOn
  isLikeOn
  ticketCount
  purchasedDate
  lastViewInfo {
    ...LastViewInfoFragment
  }
  purchaseInfo {
    ...PurchaseInfoFragment
  }
  ticketInfo {
    price
    discountPrice
    ticketType
  }
}

fragment ContinueInfoFragment on ContinueInfo {
  title
  isFree
  productId
  lastReadProductId
  scheme
  continueProductType
  hasNewSingle
  hasUnreadSingle
}

fragment TodayGift on TodayGift {
  id
  uid
  ticketType
  ticketKind
  ticketCount
  ticketExpireAt
  ticketExpiredText
  isReceived
  seriesId
}

fragment PreviewFragment on Preview {
  item {
    ...PreviewSingleFragment
  }
  nextItem {
    ...PreviewSingleFragment
  }
  usingScroll
}

fragment PreviewSingleFragment on Single {
  id
  productId
  seriesId
  title
  thumbnail
  badge
  isFree
  ageGrade
  state
  slideType
  lastReleasedDate
  size
  pageCount
  isHidden
  remainText
  isWaitfreeBlocked
  saleState
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
}

fragment OperatorPropertyFragment on OperatorProperty {
  thumbnail
  copy
  helixImpId
  isTextViewer
  selfCensorship
  isBook
  cashInfo {
    discountRate
    setDiscountRate
  }
  ticketInfo {
    price
    discountPrice
    ticketType
  }
}

fragment AssetPropertyFragment on AssetProperty {
  bannerImage
  cardImage
  cardTextImage
  cleanImage
  ipxVideo
  bannerSet {
    ...BannerSetFragment
  }
  cardSet {
    ...CardSetFragment
  }
  cardCover {
    ...CardCoverFragment
  }
}

fragment BannerSetFragment on BannerSet {
  backgroundImage
  backgroundColor
  mainImage
  titleImage
}

fragment CardSetFragment on CardSet {
  backgroundColor
  backgroundImage
}

fragment CardCoverFragment on CardCover {
  coverImg
  coverRestricted
}

fragment WaitfreeTicketFragment on WaitfreeTicket {
  chargedPeriod
  chargedCount
  chargedAt
}

fragment LastViewInfoFragment on LastViewInfo {
  isDone
  lastViewDate
  rate
  spineIndex
}

fragment PurchaseInfoFragment on PurchaseInfo {
  purchaseType
  rentExpireDate
  expired
}

fragment TranslatePropertyFragment on TranslateProperty {
  category {
    ...LocaleMapFragment
  }
  sub_category {
    ...LocaleMapFragment
  }
}

fragment LocaleMapFragment on LocaleMap {
  ko
  en
  th
}

fragment NearItemFragment on NearItem {
  productId
  slideType
  ageGrade
  isFree
  title
  thumbnail
}

fragment TextViewerData on TextViewerData {
  type
  atsServerUrl
  metaSecureUrl
  contentsList {
    chapterId
    contentId
    secureUrl
  }
}

fragment TalkViewerData on TalkViewerData {
  type
  talkDownloadData {
    dec
    host
    path
    talkViewerType
  }
}

fragment ImageViewerData on ImageViewerData {
  type
  imageDownloadData {
    ...ImageDownloadData
  }
}

fragment ImageDownloadData on ImageDownloadData {
  files {
    ...ImageDownloadFile
  }
  totalCount
  totalSize
  viewDirection
  gapBetweenImages
  readType
}

fragment ImageDownloadFile on ImageDownloadFile {
  no
  size
  secureUrl
  width
  height
}

fragment VodViewerData on VodViewerData {
  type
  vodDownloadData {
    contentPackId
    drmType
    endpointUrl
    width
    height
    duration
  }
  drmInfo {
    type
    serverType
    error
    fairplayCertificateUrl
    widevineLicenseUrl
    fairplayLicenseUrl
    token
    provider
    assertion
  }
}

fragment DisplayAd on DisplayAd {
  sectionUid
  bannerUid
  treviUid
  momentUid
}
`

// ============================================================================
// MODELS (mình chỉ khai báo những field cần dùng)
// ============================================================================

type GraphQLPayload struct {
	OperationName string                 `json:"operationName"`
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
}

type GQLResponse struct {
	Data struct {
		ViewerInfo ViewerInfo `json:"viewerInfo"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type ViewerInfo struct {
	Item       Item       `json:"item"`
	SeriesItem *Series    `json:"seriesItem,omitempty"`
	PrevItem   *NearItem  `json:"prevItem,omitempty"`
	NextItem   *NearItem  `json:"nextItem,omitempty"`
	ViewerData ViewerData `json:"viewerData"`
}

type Item struct {
	Title string `json:"title"`
	// thêm field nếu cần về sau
	ProductId int `json:"productId,omitempty"`
	SeriesId  int `json:"seriesId,omitempty"`
}

type Series struct {
	SeriesId int    `json:"seriesId"`
	Title    string `json:"title"`
}

type NearItem struct {
	ProductId int    `json:"productId"`
	Title     string `json:"title"`
	// slideType, etc. omitted — thêm nếu cần
}

type ViewerData struct {
	AtsServerUrl string        `json:"atsServerUrl"`
	MetaSecure   string        `json:"metaSecureUrl,omitempty"`
	ContentsList []ContentPart `json:"contentsList"`
}

type ContentPart struct {
	ChapterId int    `json:"chapterId,omitempty"`
	ContentId int    `json:"contentId,omitempty"`
	SecureUrl string `json:"secureUrl"`
}

type ContentFileResponse struct {
	ContentInfo struct {
		ParagraphList []ParagraphNode `json:"paragraphList"`
	} `json:"contentInfo"`
}

type ParagraphNode struct {
	Type     string          `json:"type"`
	Text     string          `json:"text"`
	Children []ParagraphNode `json:"childParagraphList,omitempty"`
}

// ============================================================================
// MAIN: loop theo nextItem.productId — giữ flow tải/chạy y nguyên code gốc
// ============================================================================

func main() {
	// context tổng cho toàn bộ chạy
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	client := &http.Client{Timeout: 30 * time.Second}
	headers := map[string]string{
		"accept":          "*/*",
		"accept-language": "vi-VN,vi;q=0.9,en-US;q=0.8,en;q=0.7",
		"content-type":    "application/json",
		"origin":          "https://page.kakao.com",
		"referer":         "https://page.kakao.com/",
		"user-agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36",
		"cookie":          CookieString,
	}

	// Bắt đầu từ đây: set seriesId và productId chương ban đầu bạn muốn crawl
	currentSeries := 61097897
	currentProduct := 64488964

	for {
		fmt.Printf("\n>>> Crawl productId = %d (seriesId=%d)\n", currentProduct, currentSeries)

		// 1) build payload giữ nguyên như bạn yêu cầu
		payload := GraphQLPayload{
			OperationName: "viewerInfo",
			Query:         ViewerQuery,
			Variables: map[string]interface{}{
				"seriesId":  currentSeries,
				"productId": currentProduct,
			},
		}
		payloadBytes, _ := json.Marshal(payload)

		respBody, err := makeRequest(ctx, client, "POST", GraphQLEndpoint, headers, bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Printf(">>> ERROR: GraphQL request failed: %v\n", err)
			break
		}

		var gqlResp GQLResponse
		if err := json.Unmarshal(respBody, &gqlResp); err != nil {
			fmt.Printf(">>> ERROR: Parse Error: %v\nBody: %s\n", err, string(respBody))
			break
		}
		if len(gqlResp.Errors) > 0 {
			fmt.Printf(">>> ERROR: GraphQL Error: %s\n", gqlResp.Errors[0].Message)
			break
		}

		viewerData := gqlResp.Data.ViewerInfo.ViewerData
		itemInfo := gqlResp.Data.ViewerInfo.Item
		nextItem := gqlResp.Data.ViewerInfo.NextItem

		if viewerData.AtsServerUrl == "" {
			fmt.Println(">>> ERROR: Không tìm thấy atsServerUrl. Kiểm tra lại Cookie.")
			break
		}

		// Tạo file
		safeTitle := sanitizeFilename(itemInfo.Title)
		fileName := safeTitle + ".txt"

		file, err := os.Create(fileName)
		if err != nil {
			fmt.Printf(">>> ERROR: Không thể tạo file %s: %v\n", fileName, err)
			break
		}
		fmt.Printf(">>> Đã tạo file: %s\n", fileName)
		fmt.Printf(">>> Tổng số phần: %d. Bắt đầu tải...\n", len(viewerData.ContentsList))

		// Ghi tiêu đề
		if _, err := file.WriteString(fmt.Sprintf("# %s\n\n", itemInfo.Title)); err != nil {
			fmt.Printf(">>> WARNING: lỗi ghi tiêu đề: %v\n", err)
		}

		// Loop tải các phần
		for i, part := range viewerData.ContentsList {
			if i > 0 {
				time.Sleep(1 * time.Second)
			}
			fullUrl := viewerData.AtsServerUrl + part.SecureUrl
			fmt.Printf(" -> [%d/%d] Đang tải... (Delay 1s)\n", i+1, len(viewerData.ContentsList))

			partBody, err := makeRequest(ctx, client, "GET", fullUrl, headers, nil)
			if err != nil {
				fmt.Printf("   [Lỗi tải phần %d]: %v\n", i+1, err)
				continue
			}

			var contentResp ContentFileResponse
			if err := json.Unmarshal(partBody, &contentResp); err != nil {
				// nếu không parse được JSON -> có thể là image / binary -> bỏ qua
				continue
			}

			for _, node := range contentResp.ContentInfo.ParagraphList {
				text := extractTextRecursive(node)
				text = strings.ReplaceAll(text, "&nbsp;", " ")
				if _, err := file.WriteString(text); err != nil {
					fmt.Printf("   Lỗi ghi file: %v\n", err)
				}
			}
		}

		file.Close()
		fmt.Println(">>> XONG! Kiểm tra file:", fileName)

		// Nếu không có nextItem -> dừng
		if nextItem == nil || nextItem.ProductId == 0 {
			fmt.Println(">>> HẾT CHƯƠNG hoặc không tìm thấy nextItem.")
			break
		}

		// Cập nhật productId để crawl chương tiếp theo
		fmt.Printf(">>> Next: %d (%s)\n", nextItem.ProductId, nextItem.Title)
		currentProduct = nextItem.ProductId

		// delay an toàn trước lần tiếp theo
		time.Sleep(2 * time.Second)
	}
}

// ============================================================================
// HELPERS
// ============================================================================

func makeRequest(ctx context.Context, client *http.Client, method, url string, headers map[string]string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Status %d | Body: %s", resp.StatusCode, string(respBytes))
	}
	return respBytes, nil
}

func extractTextRecursive(node ParagraphNode) string {
	var sb strings.Builder
	if node.Text != "" {
		sb.WriteString(node.Text)
	}
	if node.Type == "BR" {
		sb.WriteString("\n")
	}
	for _, child := range node.Children {
		sb.WriteString(extractTextRecursive(child))
	}
	switch node.Type {
	case "P", "DIV":
		sb.WriteString("\n")
	case "H3":
		sb.WriteString("\n\n### ")
	}
	return sb.String()
}

func sanitizeFilename(name string) string {
	reg, _ := regexp.Compile(`[\\/:*?"<>|]`)
	safe := reg.ReplaceAllString(name, "_")
	return strings.TrimSpace(safe)
}
