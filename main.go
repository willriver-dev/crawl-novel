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
	CookieString = "_kpdid=ca05be6623084b93b8b01c9e1200e0c3; _kpiid=8bc3c90e0a818233479d4791be94e941; _kahai=ba4be1e91c915262b60a6402f5e1e321446b1a1739f6a79ad21e84f1393c8927; _kawask=941fa7a0-bf48-4954-8eab-305130ea9bec; _kau=4d731b8667dc763029f68052ca90bf3e04c2fda5f215c9c5d508b470a073c65a9c33247f56343a1514cdc0dd6c6ce5d3ccc748eddeb90b9eba23363e673c201558cb82bc6ca281aad0a4dd4772a62067f6ec4c77266de7c7dde24aee5dfab84836690e0dfbf4b7bed6982f8fc5cb6a68f09dce60ca3a3af85aa2ba533132383432313233363436383332373835383038333133373237343233313131e536ec5f239f5be4eca10103d34d3968; _kawlt=rnkeYKoDFdqmnyPuLNkan6NIMZ37cE6OEl_40VAh3V33VFiHpidyxU9j-4eTRAvtRhZ2zVVorQfQUHoY3-uBMHhswtQTR29jVFF_58R10ZDIwyUulK2yjGiOSnqkzbhw; _kawltea=1766302547; _karmt=6f1Qj5mvMOPMtjAsm8IRrDP8VT-ZDDkcz17B46hVXmbKpvO-DFl2fA9ee0bCDglI; _karmtea=1768829747; _kaslt=W4MqDkW++Vj+s/3c6Mc3yPuChNjUuic7p7aqckCPaV7zW+Vdva0qwVeRs7J2dferiwCyUlurz1B2aiJIljYQSw==; _kpawbat_e=ZRhMkByZWy2Vh4zbshztM1MBsLcLp8CVgjVslz6IXdEym07xJxDOmoJIkxUTBksaP%2BN7JmJ5kNlNpd14wvXBK8qZdvfbLMWEV9hpHJQ6HZQ%3D; _kpwtkn=eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNIn0..AwMCKzId6kCPikzX.gLYO_Znxm-jphg-V7LC_4rPFlt18-TtkYL4orwd-12Hwb6_9_ipOTk1rxc8IM1IHP9NHezBp8WXc4PJsZbkAj5Acnf_UmW9QUoju4XHG668cZkA4Eufvcj6MSjAUr_T1wm5zXAks7_63z5dOHjWdhrA9c6MW_c4JFqIPn-5EDx9kGIfAFgINaOFCGmEZ-__EYjYe5sGjGP30VhIDsmFxLjaMLaH0BHu7fhxnLXhioV3pLuI3MQgplQeAwJ9wFfbJ70W7LVSci_pRs5ZYdp8YxO6PkVHCRkULGe9WySqTnmQ-cR-3mgaB4jq0lwqB-mkji3QrrbKvycre6RedsEbrNCUssmRvFFWh_X4xXnY0lk3q9WwJdR6hTxyqTuGDFh_0LkO-7cGQ8zLvw9BZOG3mdaleX4RK9spoPzv8W7s3cMfOvKttHfV8f-cpoXDyYrP1hM-ijimsL1x66MIkn2KVJUW97bcxQB6ECnuyUWzRHDQ3m7SI21AMM4QJRAiNx1jnFKfOtYNDeWHCJr3q7etLzT1xRjYtEb16DilRAQqyFY_JEyfuu10ARQjrIWXtTLP3XwEPRFt5bcBgVXqD1ZBQYxCPNrTAktnSD0JEpGGa3f_jfcaqyiypo7stzawQ7sk_deccChsoRd1G4VcwyYUiz9xLpjdGzWyxr3vko_CUGvTVO715bBYYRe1ZcypEFs7Esdr7sRMpA2bauWWvi_E9T0l5pcy6H0jmaOtGHagIXHzU4bW_A4LeCpmoB102PunI2anDa6afmLtJ8Uhi5xxU2Ddda6ubeDljkrWdp14sfrvHlk5USA4OmbhkbcQCkBpErgYPJdjouWMq5ImdV9YP4v4YNCmXaEhQDJsbAVzBhtmB6uM4vEf867BXpCNqR1VqItYAy5Qm84RuXyVQaudfgxw8tkBTUnbGW3HNWviK5hhZ2ClvT2i_CJyeG0wQpsPzpvQznVZvIg4Jw1ruowwm5944a8XMHkknRYDvC8LgJ4jo4sY2Yzp-tYj0bHBA86mvS-PM6b360HekTWyCSKixcNON5tqOV_-C7AaeGAbkavOBBuclwaWDo1-h5ZDdA434Zd7uc7zRGsoFzFWBCR3qsEIDSlqZklp85FcmOBcYpwGiwaZ-ra_C_9S4w1CbFZEEiHBOYQldYcFUcX6qr3Q4sjP_l5EtZb9SXv06jMva_1RBhj08YvvK4I2-XzfWrQ9A51o7RcXKZILWhz4yGs3_zZfPxWootgYLc8_WYvEYxU1GcvZpkhvLebiUwE5Zvk9u2b1HpQMY5Y3SFUz1b1TSvJ1l2qgWgwvEeKHdJGZoHHvrSewWU2CDtL7q9PpSvfEHjgzdVezIaahDXLjYL4_5K-TFkE4Hw0ln0iQe9pXcrhSeeNqRJxLA9JGt1jKD4iPDKtqCJDuG9wzLd4G1-8uAUmSSF9I_sADePEcr_a-3hQV4bVez57gxJyP3o7HhWiI4DPyh-QHJiSkJM-vF_RZPyGdBHMkP0JazVaDrUCt7bklDdgR1wl3ItsYNG9j89scBDUnnUo1YMFuFWZIxJq_POq0sg5hwta76q-qoxEci00PuHTdBWLxKD_Rz3nT6vE-0WeuU-YVzZZxyDxSbE1JUmNCwEnI15XiRrKJnK7xS5GqiS-GtR-Q141MmnoGDiszZMrHTW3uvgWUNSHGxRsKsW7mKtySIwwuIppb9vh_O9dDJjl6R24AK0AKuSH31sV28kH5Py7viQmWKgVDyH-LcpUCDOV1xv7CTSkekku3uQ30AnpLJmxO_5tkK6b6rpmYY3TFwEsfgvG5AllqFdmRQsVuHFCLDxIp0VxU9YAmHw-qHrLhSUDF9-NCHu7mI7yXJoPupvXTh_PH79J9HriSTQw1hL4w-oZLPfR7s5xv1sqvAKAC5YV1fR3sT4bENYTgTdD9CPgbAMgS2Rfi1yUrofSXs0rXJgYTbRdAhMUcKhYht8MlptuFhggXpydPDkriAPZgibT_7-Ej3FH64Fcg_qZssZ0rExp-rKTW6uYybSwkBhRgPCuRUcc_ytHTuHzv-ZC4WdeosLy_eabbgfhva7T3MmnAP7pxepcSeaOhzv85yWMfkBHoqlunAMDXXDq6i1zmIHprQHgUG0LyMYlaftYHB-aH88k9CLJuTFnHdRYEmzqsB9l4ONiBiFSGav9_g2j_WqdCmvW5UghR47-EowN2hOs7jwgz3wOb08jrOTddgBhdhVhW2mAl0VQ-IyM2B5moUC_860BVH6DsLjVDa6VKYVHXeVBnM8VDUdM4vHxfBajNvE3XEcX4XgDZQz-o.z7ZPGY1DWaqJj-j4q6p_QA; _T_ANO=US96rXkxzqFvGlBSv9z61hbM5avO6ropWUOb+XHQauplY440tB8zLL0+tlZZaFuHTHyWk/DwLayDX+9EM4m2Q84Ye0GjIBK7rXqSREb70cf14dkRaInNjLViokDSh9SU8iw18yJxu36WiT/xRC546zQ4J/5L6X2Urh8a4w3NbgvaBmZEBD7gc2TV9AHG9VhpeHXcEwT6l+/1suR5n+a02+C4qvEC0WUTLTqA3qYzhWDDyt6edrPov2FOYAw3tj1vjVDr3iFgzSDJLCPPsi2U5547DAH1ZoO7j6UI3nXYQT29UMzcIuc96I9w87fG+Zo/3nEEv60ZhITgJKhhokpXCg=="
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
	currentProduct := 65489216

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
